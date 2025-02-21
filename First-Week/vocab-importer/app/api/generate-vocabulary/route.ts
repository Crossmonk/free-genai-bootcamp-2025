import { generateText } from "ai";
import { createOpenAI as createGroq } from "@ai-sdk/openai";
import { NextResponse } from "next/server";

const groq = createGroq({
  baseURL: "https://api.groq.com/openai/v1",
  apiKey: process.env.GROQ_API_KEY,
});

export async function POST(req: Request) {
  const { category } = await req.json();

  try {
    const { text } = await generateText({
      model: groq("llama-3.3-70b-versatile"),
      prompt: `Generate a list of 5 Japanese vocabulary words related to the category "${category}". 
      For each word, provide the kanji, romaji, English translation, and its component parts (if applicable). 
      Format the response as a JSON array with the following structure:
      [{"kanji":"[japanese kanji]","romaji":"","english":"","parts":[{"kanji":"","romaji":""},{"kanji":"","romaji":""}]}]. 

      Example:
        [
          {
            "kanji": "勉強",
            "romaji": "benkyou",
            "english": "study",
            "parts": [
              { "kanji": "勉", "romaji": "ben" },
              { "kanji": "強", "romaji": "kyou" }
            ]
          },
          {
            "kanji": "食事",
            "romaji": "shokuji",
            "english": "meal",
            "parts": [
              { "kanji": "食", "romaji": "shoku" },
              { "kanji": "事", "romaji": "ji" }
            ]
          },
          {
            "kanji": "旅行",
            "romaji": "ryokou",
            "english": "travel",
            "parts": [
              { "kanji": "旅", "romaji": "ryo" },
              { "kanji": "行", "romaji": "kou" }
            ]
          },
          {
            "kanji": "運動",
            "romaji": "undou",
            "english": "exercise",
            "parts": [
              { "kanji": "運", "romaji": "un" },
              { "kanji": "動", "romaji": "dou" }
            ]
          },
          {
            "kanji": "図書",
            "romaji": "tosho",
            "english": "books",
            "parts": [
              { "kanji": "図", "romaji": "to" },
              { "kanji": "書", "romaji": "sho" }
            ]
          }
        ]
      
      Present json without any additional text or format.`,
    });

    const vocabularyList = JSON.parse(text);
    return NextResponse.json(vocabularyList);
  } catch (error) {
    console.error("Error generating vocabulary:", error);
    return NextResponse.json(
      { error: "Failed to generate vocabulary" },
      { status: 500 }
    );
  }
}
