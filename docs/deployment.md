# Compilation & Deployment Guide

This guide details how to build `jorbites-scripts` for local execution and how to cross-compile it for a remote target like a Raspberry Pi 3B.

---

## 1. Build Locally
To compile the CLI binary for your current operating system and architecture:
```bash
make build
```
This runs `go build -o jorbites-scripts main.go`. You can run the resulting executable immediately:
```bash
./jorbites-scripts list-all-badges
```

---

## 2. Cross-Compile for Raspberry Pi 3B
When compiling on your development machine (macOS, Windows, or Linux x86) to run on a Raspberry Pi 3B, you must set target environment variables (`GOOS`, `GOARCH`, and optionally `GOARM`) before running the build command.

You can use the predefined Makefile targets or run `go build` directly.

### A. For Raspberry Pi OS 32-bit (ARMv7)
If your Pi is running a 32-bit operating system:
```bash
# Using Makefile
make cross-compile

# Or manual build command
GOOS=linux GOARCH=arm GOARM=7 go build -o jorbites-scripts-pi32 main.go
```

### B. For Raspberry Pi OS 64-bit (ARM64)
If your Pi is running a 64-bit operating system:
```bash
# Using Makefile
make cross-compile

# Or manual build command
GOOS=linux GOARCH=arm64 go build -o jorbites-scripts-pi64 main.go
```

---

## 3. Deploying to the Pi
1. Copy the compiled binary (e.g. `jorbites-scripts-pi32` or `jorbites-scripts-pi64`) to the Raspberry Pi using `scp` or `sftp`:
   ```bash
   scp jorbites-scripts-pi64 pi@raspberrypi.local:/home/pi/jorbites-scripts
   ```
2. SSH into your Raspberry Pi and ensure the binary has execution permissions:
   ```bash
   chmod +x /home/pi/jorbites-scripts
   ```
3. Create a `.env` file in the same directory on the Pi containing your production connection strings:
   ```env
   DATABASE_URL="mongodb+srv://user:pass@cluster.mongodb.net/prod"
   JORBITES_URL="https://jorbites.com"
   ```
4. Run the script:
   ```bash
   ./jorbites-scripts list-badges 6fsb79dsbgsdb0bsde
   ```
