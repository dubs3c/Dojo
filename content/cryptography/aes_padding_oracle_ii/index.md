---
title: "AES Padding Oracle II"
description: "Learn how to perform the classic padding oracle attack"
date: 2020-01-27T13:33:37Z
image: hacker2.jpg
author: "dubs3c"
categories: "Cryptography"
tags: ["AES-CBC", "oracle", "padding"]
---

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam euismod et metus vitae sagittis. Fusce imperdiet molestie odio. Suspendisse at vulputate magna. Fusce semper rutrum felis, a placerat est dapibus id. Ut porta cursus dapibus. Pellentesque quis porttitor tortor. Vivamus sollicitudin purus sit amet purus pharetra ornare. Proin ipsum felis, dapibus ut egestas nec, aliquam sit amet sapien. Nam rutrum eros odio, a efficitur ligula aliquet sit amet. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae;

Phasellus ullamcorper leo et tortor placerat, eget rutrum dui congue. Nunc mi libero, interdum et erat in, tempus dignissim eros. Integer imperdiet neque non enim lacinia mollis. Vivamus pretium pulvinar sodales. Praesent suscipit quam at augue porta, a consequat justo facilisis. Mauris diam elit, tristique quis ultricies pretium, finibus eget turpis. Ut justo sapien, porta sit amet sapien nec, gravida aliquam est. Quisque fringilla nibh in nisi luctus sodales. Phasellus ullamcorper lorem eget lacus eleifend varius. Phasellus at posuere est. Aliquam interdum risus vel egestas feugiat. Nunc aliquam elementum dui, vel mollis ligula sodales in. Nullam id nunc ligula. Fusce quis velit et massa cursus interdum in in mauris. Duis ut dolor non ipsum convallis convallis. Nulla facilisi.

## Test
Nullam sed risus accumsan, hendrerit justo eu, placerat elit. Etiam porttitor mattis dui, pellentesque volutpat augue rhoncus et. Proin eget nisl euismod, efficitur tellus sit amet, facilisis ante. Donec cursus lacus at sem posuere condimentum. Maecenas hendrerit rutrum justo. Nulla facilisi. Nulla lacinia pretium erat eu sollicitudin. Nam mauris quam, laoreet vitae tristique ac, efficitur vel lacus. Aenean elementum ligula vel ligula luctus, imperdiet egestas sem semper. Ut sed metus imperdiet felis pulvinar accumsan. Proin gravida nec sem faucibus iaculis. Vestibulum feugiat dictum neque sed molestie.
### Another test
Nullam non lectus a nibh ornare luctus. Proin vitae ultricies urna, a condimentum elit. Donec sit amet suscipit mauris. Integer quis augue ut purus hendrerit tempus. Duis auctor, diam quis ornare tristique, est turpis ultricies turpis, luctus finibus nibh dui vel ex. Quisque bibendum in urna vel ornare. Duis interdum leo eleifend, congue urna ut, tempus orci. Nulla egestas aliquet arcu, sed bibendum dolor semper non. Suspendisse potenti. Sed quis lectus ut est sagittis laoreet in sed leo.

#### Another test
Nullam non lectus a nibh ornare luctus. Proin vitae ultricies urna, a condimentum elit. Donec sit amet suscipit mauris. Integer quis augue ut purus hendrerit tempus. Duis auctor, diam quis ornare tristique, est turpis ultricies turpis, luctus finibus nibh dui vel ex. Quisque bibendum in urna vel ornare. Duis interdum leo eleifend, congue urna ut, tempus orci. Nulla egestas aliquet arcu, sed bibendum dolor semper non. Suspendisse potenti. Sed quis lectus ut est sagittis laoreet in sed leo. 

```python
def classic_attack(iv: str, ciphertext_hex: str) -> bytes:
    """Performs a classic padding oracle attack where the oracle returns if padding is correct or not
    Assumes IV is known in order to recover block zero.

    Args:
        iv (bytes): Hex encoded IV. Enter "" if not known
        ciphertext_hex (str): Hex encoded ciphertext

    Returns:
        bytes: Bytes encoded plaintext
    """

    blocks = create_blocks(ciphertext_hex)
    if iv != "":
        blocks = [iv] + blocks
    intermediates = {}
    plaintext = ""

    for block_number in range(len(blocks)-2, -1, -1):

        if block_number not in intermediates:
            intermediates.update({block_number: []})

        # Always remove last block after successful decryption
        # Easier to crack remaining blocks
        mod_blocks = blocks[:block_number+2]

        # This is the current block being cracked
        work_block = bytearray(unhexlify(mod_blocks[block_number].encode()))

        for result in attack_single_block(work_block, mod_blocks, blocks, intermediates, block_number):
            plaintext += result

    return plaintext[::-1].encode()

```
