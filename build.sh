#!env sh
set +e

pushd cloud-init
CGO_ENABLED=0 go build -o ../cloudinit
popd

strip cloudinit

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

# finally include the hypervisor-fw in the image
curl -L -o output/hypervisor-fw https://github.com/cloud-hypervisor/rust-hypervisor-firmware/releases/download/0.3.2/hypervisor-fw
