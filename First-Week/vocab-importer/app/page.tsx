"use client"

import type React from "react"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

type VocabularyItem = {
  kanji: string
  romaji: string
  english: string
  parts: { kanji: string; romaji: string }[]
}

export default function VocabularyGenerator() {
  const [category, setCategory] = useState("")
  const [result, setResult] = useState<VocabularyItem[]>([])
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)

    try {
      const response = await fetch("/api/generate-vocabulary", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ category }),
      })

      if (!response.ok) {
        throw new Error("Failed to generate vocabulary")
      }

      const data = await response.json()
      setResult(data)
    } catch (error) {
      console.error("Error:", error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleCopy = () => {
    navigator.clipboard.writeText(JSON.stringify(result, null, 2))
  }

  const handleDownload = () => {
    const blob = new Blob([JSON.stringify(result, null, 2)], { type: "application/json" })
    const url = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = "vocabulary.json"
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Vocabulary Generator</h1>
      <form onSubmit={handleSubmit} className="mb-4">
        <Input
          type="text"
          value={category}
          onChange={(e) => setCategory(e.target.value)}
          placeholder="Enter thematic category"
          className="mb-2"
        />
        <Button type="submit" disabled={isLoading}>
          {isLoading ? "Generating..." : "Generate Vocabulary"}
        </Button>
      </form>
      {result.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Generated Vocabulary</CardTitle>
          </CardHeader>
          <CardContent>
            <Textarea value={JSON.stringify(result, null, 2)} readOnly className="h-64 mb-2" />
            <div className="flex gap-2">
              <Button onClick={handleCopy}>Copy to Clipboard</Button>
              <Button onClick={handleDownload}>Download JSON</Button>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}

