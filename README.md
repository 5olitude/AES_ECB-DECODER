# **AES ECB Decoder in Go**

An example of decoding an image encrypted using **AES** in **ECB** (Electronic Codebook) mode with **Go**. This demonstrates how AES-ECB works and explains why it's considered insecure for real-world usage.

## **What is AES with ECB Encryption?**

**AES** (Advanced Encryption Standard) is a symmetric block cipher that encrypts data in blocks.  
The AES block size is always **128 bits**, which is equal to **16 bytes**.

### **Process in Each Round:**

Each encryption round in AES involves the following steps:

1. **SubBytes**: A non-linear substitution step where each byte in the block is replaced using a substitution table (S-box).
2. **ShiftRows**: Rows of the state (the 16-byte block treated as a 4x4 matrix) are shifted.
3. **MixColumns**: A mixing operation combining the bytes in each column of the state (except in the final round).
4. **AddRoundKey**: The round key, derived from the encryption key, is XOR-ed with the block.

### **Ciphertext:**

After all the rounds are complete, a **16-byte (128-bit) block of ciphertext** is produced.  
If the plaintext exceeds 16 bytes, AES processes the data in **multiple 16-byte blocks**, generating 16 bytes of ciphertext for each block.

---

## **ECB (Electronic Codebook) Mode:**

ECB is one of the simplest block cipher modes of operation. In ECB mode:

- The plaintext is divided into **16-byte blocks**.
- Each block is **independently encrypted** using the same key.
- If the plaintext is not a multiple of 16 bytes, padding is applied to the final block.
- The resulting ciphertext blocks are concatenated to form the final encrypted message.

### **How AES-ECB Works:**

1. **Plaintext** is divided into 16-byte blocks.  
2. **Each block** is independently encrypted using AES and the same key.
3. The **resulting ciphertext** blocks are concatenated.

   ![Wikipedia ECB Penguin Image](https://github.com/5olitude/AES_ECB-DECODER/blob/7a2a45ad849d19343d4fe402bba6ac78275a88f3/Screens.png)
---

## **Real-World Example: ECB Penguin Attack**

The **ECB Penguin** is a famous example illustrating the weakness of AES-ECB mode. When an image (like the famous penguin image) is encrypted using ECB, **repeating patterns** in the plaintext are easily visible in the ciphertext. This happens because identical plaintext blocks are encrypted into identical ciphertext blocks, making the encrypted image resemble the original.  
Learn more about the **ECB Penguin** from this [GitHub link](https://github.com/robertdavidgraham/ecb-penguin).

You can also refer to the following Wikipedia illustration:



### **Other Real-World Attacks**:

One example is the **Adobe Password Database Leak**, where identical passwords produced identical ciphertexts. This allowed attackers to locate frequently reused passwords, leading to significant data breaches.

Here is a comic from **XKCD** that humorously illustrates this vulnerability:

![XKCD Encryptic Comic](https://imgs.xkcd.com/comics/encryptic_2x.png)

For a detailed explanation, refer to this [StackExchange discussion](https://crypto.stackexchange.com/questions/14487/can-someone-explain-the-ecb-penguin).
