# Projget

Projget is a lightweight project bundling and retrieval tool written in Go.  
It is designed as a simplified alternative to Git for developers who want a minimal workflow focused on packaging, fetching, and distributing project directories.

Projget does not implement full version control. Instead, it provides a straightforward way to initialize a project, bundle it into a portable format, retrieve remote bundles, and extract them.

---

## Features

- Simple project initialization
- Custom bundle format (.pf)
- Magic identifier validation
- Fetch remote bundle files over HTTP/HTTPS
- Extract bundled projects
- Minimal dependencies
- Single-binary distribution

---

## Project Structure

When initialized, a project contains:
Dir/
 files <like main.c hi.go>
.Proj/
  indent.txy

  The `.Proj/ident.txt` file contains a magic number in hexadecimal format to validate the project structure.

Bundled projects are packaged as: Dir.pf

---

## Installation

Make sure you have Go installed (1.20+ recommended).

Build from source:

```bash
go build -o projget main.go
```
This produces a single executable named `projget`.

---

## Commands

### init

Initializes a project in the current directory.

projget init

Creates:

- `.Proj/`
- `.Proj/ident.txt` containing the magic number

---

### bundle

Bundles the current directory into a `.pf` file.
projget bundle

Output:
<directory_name>.pf

---

### get

Downloads a remote `.pf` bundle via HTTP or HTTPS.

projget get 

Example:

projget get https://example.com/project.pf�

---

### unbundle

Extracts a `.pf` bundle into a directory.

projget unbundle <file.pf>

---

### push

Uploads a bundled directory to a remote endpoint (if supported).

projget push 

---

## Bundle Format

Projget bundles are stored with:

- A magic header identifier
- Directory structure metadata
- File contents

The magic identifier ensures that only valid `.pf` bundles are extracted.

---

## Use Cases

- Lightweight project sharing
- Simple distribution of source code
- Educational tooling
- Alternative workflow for small projects

---

## Limitations

- No commit history
- No branching
- No merge support
- No diff comparison
- No conflict resolution

Projget is not a replacement for Git in large collaborative environments. It is intended for simplicity and portability.

---

## Roadmap

Planned improvements may include:

- File hashing
- Incremental bundling
- Project status command
- Version tracking
- Improved remote push support

---

## License

MIT Lisence


 
