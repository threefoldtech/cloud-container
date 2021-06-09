root=$1
init=$2

sudo virtiofsd --socket-path=/tmp/root.socket -o source=${root} -o cache=none &

if [ -z "${init}" ]; then
    init='init=/bin/sh'
fi

exec sudo cloud-hypervisor \
    --kernel output/vmlinuz-linux \
    --initramfs output/initramfs-linux.img \
    --console off \
    --serial tty \
    --cpus boot=1 \
    --memory size=1024M,shared=on \
    --fs tag=/dev/root,socket=/tmp/root.socket  \
    --net tap=cont0,ip=192.168.123.100,mask=255.255.255.0 \
    --cmdline "console=ttyS0 rootfstype=virtiofs root=/dev/root rw ${init}" \
    --rng
