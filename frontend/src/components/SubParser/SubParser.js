import assParse from "./AssParser";
export default async function subParse(fileUrl, type) {
  if (type === ".ass") {
    try {
      let ret = await assParse(fileUrl);
      return ret;
    } catch (error) {
      return [];
    }
  } else {
    throw new Error("Unsupported subtitle file type: " + type);
  }
}