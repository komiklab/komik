const encoder = new TextEncoder();

export async function signHookRequest(
  secret: string,
  body: string,
): Promise<{
  "x-signature": string;
  "x-timestamp": string;
}> {
  const timestamp = Date.now().toString();
  const key = await crypto.subtle.importKey(
    "raw",
    encoder.encode(secret),
    { name: "HMAC", hash: "SHA-256" },
    false,
    ["sign"],
  );
  const payload = new Uint8Array([
    ...encoder.encode(timestamp),
    ...encoder.encode(body),
  ]);
  const signature = await crypto.subtle.sign("HMAC", key, payload);
  const signatureHex = Array.from(new Uint8Array(signature), (b) =>
    b.toString(16).padStart(2, "0"),
  ).join("");
  return {
    "x-signature": signatureHex,
    "x-timestamp": timestamp.toString(),
  };
}
