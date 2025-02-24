package loaders

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/richardlehane/mscfb"
	"github.com/tealeg/xlsx"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

type Office struct {
	reader   io.ReaderAt
	size     int64
	fileType string
}

var _ documentloaders.Loader = Office{}

func NewOffice(reader io.ReaderAt, size int64, filename string) Office {
	return Office{
		reader:   reader,
		size:     size,
		fileType: strings.ToLower(filepath.Ext(filename)),
	}
}

func (loader Office) Load(ctx context.Context) ([]schema.Document, error) {

	switch loader.fileType {
	case ".doc", ".docx":
		return loader.loadWord()
	case ".xls", ".xlsx":
		return loader.loadExcel()
	case ".ppt", ".pptx":
		return loader.loadPowerPoint()
	default:
		return nil, fmt.Errorf("unsupported file type: %s", loader.fileType)
	}
}

func (loader Office) LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]schema.Document, error) {
	docs, err := loader.Load(ctx)
	if err != nil {
		return nil, err
	}

	return textsplitter.SplitDocuments(splitter, docs)
}

func (loader Office) loadWord() ([]schema.Document, error) {
	if loader.fileType == ".docx" {
		return loader.loadDocx()
	}

	return loader.loadDoc()
}

func (loader Office) loadDoc() ([]schema.Document, error) {
	doc, err := mscfb.New(io.NewSectionReader(loader.reader, 0, loader.size))
	if err != nil {
		return nil, fmt.Errorf("failed to read DOC file: %w", err)
	}

	var text strings.Builder
	for entry, err := doc.Next(); err == nil; entry, err = doc.Next() {
		if entry.Name == "WordDocument" {
			buf := make([]byte, entry.Size)
			i, err := doc.Read(buf)
			if err != nil {
				return nil, fmt.Errorf("error reading WordDocument stream: %w", err)
			}
			if i > 0 {
				// Process the binary content
				for j := 0; j < i; j++ {
					// Extract readable ASCII text
					if buf[j] >= 32 && buf[j] <= 126 {
						text.WriteByte(buf[j])
					} else if buf[j] == 13 || buf[j] == 10 {
						text.WriteByte('\n')
					}
				}
			}
		}
	}

	return []schema.Document{
		{
			PageContent: text.String(),
			Metadata: map[string]interface{}{
				"fileType": loader.fileType,
			},
		},
	}, nil
}

func (loader Office) loadDocx() ([]schema.Document, error) {
	// First read into buffer
	buf := bytes.NewBuffer(make([]byte, 0, loader.size))
	if _, err := io.Copy(buf, io.NewSectionReader(loader.reader, 0, loader.size)); err != nil {
		return nil, fmt.Errorf("failed to copy content: %w", err)
	}

	// Create zip reader from buffer
	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), loader.size)
	if err != nil {
		return nil, fmt.Errorf("failed to read DOCX file as ZIP: %w", err)
	}

	var text strings.Builder
	for _, file := range zipReader.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("error opening document.xml: %w", err)
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("error reading content: %w", err)
			}

			content = bytes.ReplaceAll(content, []byte("<"), []byte(" <"))
			content = bytes.ReplaceAll(content, []byte(">"), []byte("> "))
			text.Write(content)
		}
	}

	return []schema.Document{
		{
			PageContent: text.String(),
			Metadata: map[string]interface{}{
				"fileType": loader.fileType,
			},
		},
	}, nil
}

func (loader Office) loadExcel() ([]schema.Document, error) {
	// Create a temporary buffer to store the Excel file
	buf := bytes.NewBuffer(make([]byte, 0, loader.size))
	if _, err := io.Copy(buf, io.NewSectionReader(loader.reader, 0, loader.size)); err != nil {
		return nil, fmt.Errorf("failed to copy Excel content: %w", err)
	}

	// Use tealeg/xlsx to read the Excel file
	xlFile, err := xlsx.OpenBinary(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to read Excel file: %w", err)
	}

	var docs []schema.Document
	for i, sheet := range xlFile.Sheets {
		var text strings.Builder
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text.WriteString(cell.String() + "\t")
			}
			text.WriteString("\n")
		}

		docs = append(docs, schema.Document{
			PageContent: text.String(),
			Metadata: map[string]interface{}{
				"fileType":   loader.fileType,
				"sheetName":  sheet.Name,
				"sheetIndex": i,
			},
		})
	}

	return docs, nil
}

func (loader Office) loadPowerPoint() ([]schema.Document, error) {
	if loader.fileType == ".pptx" {
		return loader.loadPptx()
	}
	return loader.loadPpt()
}

func (loader Office) loadPpt() ([]schema.Document, error) {
	doc, err := mscfb.New(io.NewSectionReader(loader.reader, 0, loader.size))
	if err != nil {
		return nil, fmt.Errorf("failed to read PPT file: %w", err)
	}

	var text strings.Builder
	for entry, err := doc.Next(); err == nil; entry, err = doc.Next() {
		// PPT text is typically in streams named "PowerPoint Document"
		if entry.Name == "PowerPoint Document" {
			buf := make([]byte, entry.Size)
			i, err := doc.Read(buf)
			if err != nil {
				return nil, fmt.Errorf("error reading PowerPoint stream: %w", err)
			}
			if i > 0 {
				// Process the binary content
				for j := 0; j < i; j++ {
					// Extract readable ASCII text
					if buf[j] >= 32 && buf[j] <= 126 {
						text.WriteByte(buf[j])
					} else if buf[j] == 13 || buf[j] == 10 {
						text.WriteByte('\n')
					}
				}
			}
		}
	}

	return []schema.Document{
		{
			PageContent: text.String(),
			Metadata: map[string]interface{}{
				"fileType": loader.fileType,
			},
		},
	}, nil
}

func (loader Office) loadPptx() ([]schema.Document, error) {
	buf := bytes.NewBuffer(make([]byte, 0, loader.size))
	if _, err := io.Copy(buf, io.NewSectionReader(loader.reader, 0, loader.size)); err != nil {
		return nil, fmt.Errorf("failed to copy content: %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), loader.size)
	if err != nil {
		return nil, fmt.Errorf("failed to read PPTX file as ZIP: %w", err)
	}

	var text strings.Builder
	for _, file := range zipReader.File {
		// PPTX stores slide content in ppt/slides/slide*.xml files
		if strings.HasPrefix(file.Name, "ppt/slides/slide") && strings.HasSuffix(file.Name, ".xml") {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("error opening slide XML: %w", err)
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("error reading content: %w", err)
			}

			content = bytes.ReplaceAll(content, []byte("<"), []byte(" <"))
			content = bytes.ReplaceAll(content, []byte(">"), []byte("> "))
			text.Write(content)
			text.WriteString("\n--- Next Slide ---\n")
		}
	}

	return []schema.Document{
		{
			PageContent: text.String(),
			Metadata: map[string]interface{}{
				"fileType": loader.fileType,
			},
		},
	}, nil
}
