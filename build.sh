#!env sh
set +e

docker build -t kernel .

tmp=$(mktemp -d)
cleanup() {
    rm -rf ${tmp}
}

trap cleanup EXIT
docker save kernel | tar -x -C ${tmp}
mkdir output

for layer in $(find ${tmp} -name layer.tar); do
    tar -xf $layer -C output
done
