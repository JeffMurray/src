package safe_file_relay implements the WriteAt and Read interfaces

It is designed to concurently connect non-sequential WrtieAts from a sftp server to sequential Reads from an s3 bucket uploader.

Each file transfer gets its own instance that insures any temporary buffer caching to disk is encrypted and cleaned up when done and Reads happen as any  missing chunks come in.