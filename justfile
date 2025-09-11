export CGO_ENABLED := "0"

binary_name := "passgen"
version := `git describe --tags --always`
commit := `git rev-parse HEAD`
build_date := datetime("%Y-%m-%d")
output_path := "dist"

default:
    @just --list

prepare:
    go mod tidy
    mkdir -p {{output_path}}

release: prepare clean
    #!/usr/bin/env bash
    set -euxo pipefail
    
    platforms=(
        "linux/amd64"
        "linux/arm64"
        "windows/amd64"
        "windows/arm64"
        "darwin/amd64"
        "darwin/arm64"
    )
    
    for platform in "${platforms[@]}"; do
        GOOS=$(echo $platform | cut -d/ -f1)
        GOARCH=$(echo $platform | cut -d/ -f2)
        SUFFIX=""
        ARCH_NAME=$GOARCH
        OUTPUT_FOLDER="{{output_path}}/$GOOS-$ARCH_NAME"
        
        [[ $GOOS == "windows" ]] && SUFFIX=".exe"
        [[ $GOARCH == "amd64" ]] && ARCH_NAME="x86_64"
        
        BINARY_NAME="{{binary_name}}$SUFFIX"
        
        mkdir -p $OUTPUT_FOLDER
        
        BINARY_PATH="$OUTPUT_FOLDER/{{binary_name}}$SUFFIX"
        echo "Building $BINARY_PATH"
        
        GOOS=$GOOS GOARCH=$GOARCH go build -trimpath \
            -ldflags="-s -w -X main.version={{version}} -X main.commit={{commit}} -X main.date={{build_date}}" \
            -gcflags="all=-l -B -wb=false" \
            -o $BINARY_PATH ./cmd/main.go
        
        if [[ $GOOS == "windows" ]]; then
            zip -j -9 "{{output_path}}/{{binary_name}}-$GOOS-$ARCH_NAME.zip" "$OUTPUT_FOLDER/$BINARY_NAME" README.md
        else
            GZIP=-9 tar -cvzf "{{output_path}}/{{binary_name}}-$GOOS-$ARCH_NAME.tar.gz" --transform='s|.*/||' "$OUTPUT_FOLDER/$BINARY_NAME" README.md 
        fi
    done
    
    (sha256sum {{output_path}}/*.tar.gz {{output_path}}/*.zip > {{output_path}}/checksums.txt)

clean:
    rm -rf {{output_path}}
