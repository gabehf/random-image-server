# random-image-server
Serves a random image from a directory. Listens to file system events to automatically update available images.

NOTE: I threw together this program to serve random backgrounds to my homelab [homepage](https://gethomepage.dev/). I have not done any testing. Software provided as-is.

## Usage
Simply run the program with no arguments.

Two environment variables are available:
- `IMAGE_DIR`: The directory to scan for images. Default: `/images`
- `ALLOWED_EXTENSIONS`: Extensions to be considered by the program. Comma-separated list. Default: `.png,.jpg,.jpeg,.webp`

## Installation
### Docker
Docker run example:
```bash
docker run -p 8080:8080 --volume /path/to/images:/images:ro gabehf/random-image-server
```

Docker compose example:
```yaml
services:
  random-image-server:
    image: gabehf/random-image-server:latest
    container_name: random-image-server
    environment:
      ALLOWED_EXTENSIONS: .png,.jpg # comma-separated; default: .png,.jpg,.jpeg,.webp
    ports:
      - 8888:8080
    volumes:
      - /path/to/my/images/:/images:ro
    restart: unless-stopped

```

### From Source
Clone the repository:
```
git clone github.com/gabehf/random-image-server
```

Download dependencies:
```
go mod download
```

Build & run the program
```
go build -o random-image-server . && ./random-image-server
```

Then navigate to `localhost:8080` and you will be served a random image from the directory.