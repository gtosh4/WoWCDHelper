declare module 'json-url' {
  interface Codec {
    compress: (obj: string) => Promise<string>
    decompress: (str: string) => Promise<string>
    stats: (
      obj: object
    ) => Promise<{ rawencoded: any; compressedencoded: any; compression: any }>
  }

  const jsonurl: (codecName: string) => Codec

  export default jsonurl
}
