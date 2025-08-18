package backend

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func ParseMarkdownNote(path, name string) (string, string, []string, time.Time, error) {
	file, err := os.Open(filepath.Join(path, name))
	if err != nil {
		return "", "", nil, time.Time{}, err
	}
	defer file.Close()

	var (
		icon    string
		tags    []string
		title   string
		created time.Time
	)

	scanner := bufio.NewScanner(file)
	var metadata strings.Builder
	inMetadata := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "+++" {
			if !inMetadata {
				inMetadata = true
				continue
			} else {
				inMetadata = false
				continue
			}
		}

		if inMetadata {
			metadata.WriteString(line + "\n")
			continue
		}

		if strings.HasPrefix(line, "# ") && title == "" {
			title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", nil, time.Time{}, err
	}

	metaDataMap := make(map[string]any)
	if err := toml.Unmarshal([]byte(metadata.String()), &metaDataMap); err != nil {
		return "", "", nil, time.Time{}, err
	}

	if val, ok := metaDataMap["icon"].(string); ok {
		icon = val
	}
	if val, ok := metaDataMap["tags"].([]any); ok {
		for _, tag := range val {
			if tagStr, ok := tag.(string); ok {
				tags = append(tags, tagStr)
			}
		}
	}
	if val, ok := metaDataMap["created"].(time.Time); ok {
		created = val
	}

	if icon == "" {
		icon = "󰎞"
	}

	if title == "" {
		title = strings.TrimSuffix(name, filepath.Ext(name))
	}

	return title, icon, tags, created, nil
}

func ReadMarkdownNote(path, name string) (string, error) {
	file, err := os.Open(filepath.Join(path, name))
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	content := string(data)

	start := strings.Index(content, "+++")
	if start != -1 {
		end := strings.Index(content[start+3:], "+++")
		if end != -1 {
			content = content[:start] + content[start+3+end+3:]
		}
	}

	return strings.TrimSpace(content), nil
}

func AddTitleMarkdownNote(path, name, title string) error {
	file, err := os.OpenFile(filepath.Join(path, name), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "\n# %s\n", title); err != nil {
		return fmt.Errorf("could not write title to markdown note: %v", err)
	}
	return nil
}

func AddMetadataMarkdownNote(path, name, metadata string, value any) error {
	file, err := os.OpenFile(filepath.Join(path, name), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	content := make([]byte, stat.Size())
	_, err = file.Read(content)
	if err != nil && err != io.EOF {
		return err
	}

	parts := strings.SplitN(string(content), "+++", 3)
	var metadataBlock, body string
	if len(parts) == 3 {
		metadataBlock = strings.TrimSpace(parts[1])
		body = strings.TrimSpace(parts[2])
	} else {
		body = string(content)
	}

	metadataMap := make(map[string]string)
	if metadataBlock != "" {
		for line := range strings.SplitSeq(metadataBlock, "\n") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				metadataMap[key] = val
			}
		}
	}

	metadataMap[metadata] = fmt.Sprintf("%v", value)

	var newMetadataBlock strings.Builder
	newMetadataBlock.WriteString("+++\n")
	for key, val := range metadataMap {
		newMetadataBlock.WriteString(fmt.Sprintf("%s = %s\n", key, val))
	}
	newMetadataBlock.WriteString("+++\n\n")

	newContent := newMetadataBlock.String() + body

	err = file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(newContent))
	if err != nil {
		return err
	}

	return nil
}
