#!/usr/bin/env python3
import re
import sys

def colorize_and_format_size(size):
    # Add comma separators for thousands
    size_with_commas = "{:,}".format(int(size))
    # Count the number of commas added
    comma_count = size_with_commas.count(',')
    # Determine the number of leading spaces needed for alignment
    if comma_count == 0:
        aligned_size = '  ' + size_with_commas
    elif comma_count == 1:
        aligned_size = ' ' + size_with_commas
    else:
        aligned_size = size_with_commas
    # Colorize the size using ANSI escape codes (e.g., blue text)
    colored_size = f"\033[34m{aligned_size}\033[0m"
    return colored_size

def colorize_filename(filename):
    # Define color codes for different file types
    colors = {
    # 91-96 for bright colors 
        '.tiff': '\033[92m',    # Slightly-brighter Green
        '/': '\033[91m',    # Slightly-brighter Red
        '.jpeg': '\033[32m',    # Green 
        '.jpg': '\033[35m',    # Magenta
        '.txt': '\033[36m',    # Cyan

        '.go': '\033[32m',    # Green
        '.py': '\033[34m',    # Purple
        '.sh': '\033[33m',    # Gold
        '.pages': '\033[96m'  # Red
    }
    # Default color (reset)
    default_color = '\033[0m'
    # Find the suffix and apply the corresponding color
    for suffix, color in colors.items():
        if filename.endswith(suffix):
            return f"{color}{filename}{default_color}"
    return filename

def process_lines(lines):
    # Process each line and apply the formatting
    for line in lines:
        match = re.match(r'(\s*\d+\s+\S+\s+\d+\s+\S+\s+)(\d+)(\s+)(.*)', line)
        if match:
            prefix, size, separator, filename = match.groups()
            formatted_size = colorize_and_format_size(size)
            colored_filename = colorize_filename(filename)
            yield f"{prefix}{formatted_size}{separator}{colored_filename}"
        else:
            yield line

def main():
    lines = [line.rstrip() for line in sys.stdin]
    for processed_line in process_lines(lines):
        print(processed_line)

if __name__ == "__main__":
    main()
