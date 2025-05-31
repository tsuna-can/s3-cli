# s3-cli

An interactive command line interface for browsing and downloading files from AWS S3 buckets.

## Features

- Browse S3 buckets and objects with an intuitive terminal UI
- Filter buckets and objects by name (case-insensitive, partial match)
- Download S3 objects directly to your local filesystem
- Support for AWS profiles
- Compatible with LocalStack for development and testing

## Installation

```bash
# Clone the repository
git clone https://github.com/tsuna-can/s3-cli.git
cd s3-cli

# Build the application
go build -o s3-cli .
```

## Usage

```bash
# Basic usage (uses default AWS profile)
./s3-cli

# Specify output directory for downloads
./s3-cli --output-dir ~/Downloads

# Use a specific AWS profile
./s3-cli --profile your-profile

# Enable debug mode
./s3-cli --debug
```

## Navigation Controls

- **↑/↓**: Navigate through buckets and objects
- **Enter**: Select a bucket or download an object
- **Esc**: Return to bucket list from object view
- **Type text**: Filter buckets or objects by name
- **Ctrl+C**: Exit the application

## Roadmap

Currently, s3-cli only supports downloading objects from S3. The following features are planned for future development:

- **Upload functionality**: Upload files to S3 buckets
- **Delete functionality**: Remove objects from S3
