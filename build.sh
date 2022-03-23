#!env sh
set +e

pushd cloud-init
CGO_ENABLED=0 go build -o ../cloudinit
popd

docker build -t kernel .

tmp=$(mktemp -d)
cleanup() {
    rm -rf ${tmp}
}

trap cleanup EXIT
docker save kernel | tar -x -C ${tmp}
mkdir -p output || true

for layer in $(find ${tmp} -name layer.tar); do
    tar -xf $layer -C output
done
