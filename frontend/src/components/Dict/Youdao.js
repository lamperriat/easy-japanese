import { API_BASE_URL } from '../../services/api';

const YOUDAO_BASE = "https://dict.youdao.com/result?word={}&lang=ja"

export async function GetYoudao(word, token) {
  const url = YOUDAO_BASE.replace("{}", encodeURIComponent(word));
  try {
    const response = await fetch(`${API_BASE_URL}/api/proxy`, {
      method: 'POST',
      headers: {
          'Content-Type': 'application/json',
          "Authorization": token,
      },
      body: JSON.stringify({
          url: url,
      })
    })
    const html = await response.text();
    const container = document.createElement('div');
    container.style.display = 'none';
    container.innerHTML = html;
    const target = container.querySelector(".search_result-dict")
      .querySelector(".modules");
    // what we need: 
    // 1. extract all elements of `each-sense` class as an array
    const senses = target.querySelector(".each-page")
      .querySelectorAll(".each-sense");
    // 2. for each of them, keep text only (as string)
    const senseTexts = Array.from(senses).map(sense => {
      return sense.innerText.trim();
    });
    // 3. join them with a line break
    let resultText = senseTexts.join('<br>') + '<br><b>例句</b><br> ';
    const example_container = target.querySelector(".blng_sents_part");
    const examples_cn = example_container.querySelectorAll(".sen-eng");
    const examples_jp = example_container.querySelectorAll(".sen-ch");
    for (let i = 0; i < Math.min(examples_cn.length, 2); i++) {
      const cn = examples_cn[i].innerText.trim();
      const jp = examples_jp[i].innerText.trim();
      resultText += `<span class="example-cn">${cn}</span><br>`;
      resultText += `<span class="example-jp">${jp}</span><br>`;
    }
    // 4. return the result as a string
    return resultText;

  } catch (error) {
    console.error("Error fetching data from Youdao:", error);
    return null;
  }
}
