---
title: "AES Basics"
description: "Learn how to perform the classic padding oracle attack"
date: 2020-01-27T13:33:37Z
image: hacker2.jpg
author: "dubs3c"
categories: "Cryptography"
tags: ["AES-CBC", "AES-ECB"]
---

## AES-ECB

Electronic Code Book (ECB) is a simple block cipher. Given a plaintext block `P` it will produce a ciphertext block `C`. Meaning, if you have multiple identical plaintext blocks, AES-ECB will produce identical ciphertext blocks. Because of this property, you don't want to encrypt information longer than `BLOCK_SIZE` because any repeating blocks will produce identical ciphertext blocks. This leaks information about the plaintext which is considered insecure. 

**Encryption**
![images/CBC_encryption.png](images/ECB_encryption.png)

**Decryption**
![images/CBC_decryption.png](images/ECB_decryption.png)

## AES-CBC

Cipher block chaining (CBC) mode is one of the more common modes used in AES. It is built upon ECB but with a twist: the previous block is XORed with the current block. The result is that ciphertexts encrypted using different IVs will have different outputs, unlike ECB mode where each plaintext block `P` will always produce the same ciphertext block `C`. This property only holds if each ciphertext encrypted with the same key has a different IV. Reusing IV for every ciphertext encrypted with key `K` will reduce CBC to ECB mode. Therefore if you are encrypting multiple blocks using CBC, make sure to have unique IVs.

**Encryption**
![images/CBC_encryption.png](images/CBC_encryption.png)

**Decryption**
![images/CBC_decryption.png](images/CBC_decryption.png)


## AES-CTR

## Exercises