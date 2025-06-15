function sureIsJapanese(text) {
  // Regular expression covering:
  // - Hiragana: \u3040-\u309F
  // - Full-width Katakana: \u30A0-\u30FF
  // - Katakana Phonetic Extensions: \u31F0-\u31FF
  // - Half-width Katakana: \uFF65-\uFF9F
  const japaneseRegex = /[\u3040-\u309F\u30A0-\u30FF\u31F0-\u31FF\uFF65-\uFF9F]/;
  return japaneseRegex.test(text);
}

function timeStrToFloat(timeStr) {
  const [hours, minutes, rest] = timeStr.split(':');
  const [seconds, hundredths] = rest.split('.');
  return (
      parseInt(hours, 10) * 3600 +
      parseInt(minutes, 10) * 60 +
      parseInt(seconds, 10) +
      parseInt(hundredths, 10) / 100
  );
}

function floatToTimeStr(timeFloat) {
  if (timeFloat < 0) {
      timeFloat = 0;
  }
  let totalSecondsInt = Math.floor(timeFloat);
  const fractional = timeFloat - totalSecondsInt;
  let hundredths = Math.round(fractional * 100);
  
  if (hundredths >= 100) {
      hundredths = 0;
      totalSecondsInt += 1;
  }
  
  const hours = Math.floor(totalSecondsInt / 3600);
  const remainingSeconds = totalSecondsInt % 3600;
  const minutes = Math.floor(remainingSeconds / 60);
  const seconds = remainingSeconds % 60;
  
  const hoursStr = hours.toString().padStart(2, '0');
  const minutesStr = minutes.toString().padStart(2, '0');
  const secondsStr = seconds.toString().padStart(2, '0');
  const hundredthsStr = hundredths.toString().padStart(2, '0');
  
  return `${hoursStr}:${minutesStr}:${secondsStr}.${hundredthsStr}`;
}

class AssLine {
  constructor(start, end, cn_text, jp_text) {
    this.start = start; 
    this.end = end; 
    this.cn_text = cn_text;
    this.jp_text = jp_text;
  }

}

export default async function assParse(fileUrl) {
  // Idea:
  // 1. filter out lines that starts with `Dialogue:`
  // format:
  // `Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text`
  // 2. split by `,` and extract: Start, End, Text
  // 3. Each interval [Start, End] corresponds to 1 line of Chinese and 1 line of Japanese
  //    Insert them into a map.
  // 4. Return an array of AssLine objects.

  let mp = new Map();
  let data;
  try {
    console.log(fileUrl);
    const response = await fetch(fileUrl);
    data = await response.text();
  } catch (error) {
    console.error('Error fetching or reading the file:', error);
    return [];
  } finally {
    console.log('destroying fileUrl:', fileUrl);
    URL.revokeObjectURL(fileUrl);
  }

  data.split('\n')
    .forEach(line => {
      if (line.startsWith('Dialogue:')) {
        const parts = line.split(',');
        const start_str = parts[1].trim();
        const end_str = parts[2].trim();
        const text = parts.slice(9).join(',').trim(); // Join the rest as text
        // text may have content `{some useless stuff...}`, clear it
        const cleanedText = text.replace(/{([^{}]*|{[^{}]*})*}/g, '').trim();
        if (cleanedText === '') {
          return; // skip empty text
        }
        
        if (!mp.has(start_str)) {
          mp.set(start_str, new AssLine(
            timeStrToFloat(start_str),
            timeStrToFloat(end_str),
            '', '')
          );
        }
        if (sureIsJapanese(cleanedText) || mp.get(start_str).cn_text !== '') {
          mp.get(start_str).jp_text = cleanedText;
        } else {
          mp.get(start_str).cn_text = cleanedText;
        }
      }
    });

  return Array.from(
    mp.values()
  ).filter(line => line.cn_text !== '' || line.jp_text !== '')
  .sort((a, b) => a.start - b.start);

  
}