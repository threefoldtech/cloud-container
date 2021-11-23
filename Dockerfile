FROM archlinux as builder
RUN pacman -Sy
RUN pacman -S --noconfirm linux mkinitcpio inetutils base-devel bc python3 pahole
WORKDIR /opt
RUN curl -O https://mirrors.edge.kernel.org/pub/linux/kernel/v5.x/linux-5.12.9.tar.gz
RUN tar -xf linux-5.12.9.tar.gz
COPY config /opt/linux-5.12.9/.config
WORKDIR /opt/linux-5.12.9/
RUN make -j $(nproc)
# this is all done later so build goes faster
# if init files has changed, since it's rarely when
# linux build is gonna change

COPY mkinitcpio.conf /root/
COPY initcpio /root/initcpio
COPY setupnetwork /
# override the original initcpio
COPY init /usr/lib/initcpio
WORKDIR /root
RUN KERNELVERSION=$(ls /lib/modules) mkinitcpio -D /usr/lib/initcpio -D initcpio -v -c mkinitcpio.conf -g initramfs-linux.img

FROM scratch
COPY --from=builder /root/initramfs-linux.img /
COPY --from=builder /opt/linux-5.12.9/arch/x86/boot/compressed/vmlinux.bin /kernel
