package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// INI represents the structure of an INI file.
// It uses a map of sections, where each section is a map of key-value pairs.
type INI struct {
	data map[string]map[string]string
	mu   sync.RWMutex // For thread-safe access
}

// NewINI creates and initializes a new INI structure.
func NewINI() *INI {
	return &INI{
		data: make(map[string]map[string]string),
	}
}

// LoadFile parses an INI file from the given path.
func LoadConf(filename string) (*INI, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	ini := NewINI()
	scanner := bufio.NewScanner(file)
	currentSection := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
			// Ensure the section map is created
			if _, ok := ini.data[currentSection]; !ok {
				ini.data[currentSection] = make(map[string]string)
			}
			continue
		}

		// Check for key-value pair
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if currentSection != "" {
				ini.data[currentSection][key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return ini, nil
}

// Get retrieves a value for a given section and key.
// It returns the value and a boolean indicating if the key was found.
func (i *INI) GetConf(section, key string) (string, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	sec, ok := i.data[section]
	if !ok {
		return "", false
	}
	val, ok := sec[key]
	return val, ok
}

// Set adds or updates a value for a given section and key.
// If the section does not exist, it will be created.
func (i *INI) SetConf(section, key, value string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.data[section]; !ok {
		i.data[section] = make(map[string]string)
	}
	i.data[section][key] = value
}

// SaveFile writes the INI data to a file at the given path.
func (i *INI) SaveConf(filename string) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for section, kv := range i.data {
		// Write section header
		if _, err := fmt.Fprintf(writer, "[%s]\n", section); err != nil {
			return err
		}
		// Write key-value pairs
		for key, value := range kv {
			if _, err := fmt.Fprintf(writer, "%s=%s\n", key, value); err != nil {
				return err
			}
		}
		// Add a blank line for readability
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}

	return writer.Flush()
}
