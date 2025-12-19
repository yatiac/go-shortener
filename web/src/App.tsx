import { useState } from 'react'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { AlertCircle, Copy, Link } from "lucide-react"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { FaGithub } from 'react-icons/fa'

export default function App() {
  const [url, setUrl] = useState('')
  const [error, setError] = useState('')
  const [shortUrl, setShortUrl] = useState('')
  const [loading, setLoading] = useState(false)
  const [copySuccess, setCopySuccess] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setShortUrl('')
    setLoading(true)

    if (!isValidUrl(url)) {
      setError('Please enter a valid URL')
      setLoading(false)
      return
    }

    try {
      const result = await fetch('/api/shorten', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ long_url: url }),
      }).then(res => res.json()).then(data => data.slug)
      setShortUrl(`${import.meta.env.VITE_BASE_URL}/${result}`)
      setUrl('') // Clear the input after successful submission
    } catch (e) {
      console.log(e)
      setError('An error occurred while shortening the URL')
    } finally {
      setLoading(false)
    }
  }

  const isValidUrl = (string: string) => {
    try {
      new URL(string)
      return true
    } catch (e) {
      if (e instanceof TypeError) {
        setError('The URL format is incorrect. Please ensure it starts with http:// or https://')
      } else {
        setError('An unexpected error occurred while validating the URL')
      }
      console.log(e)
      return false
    }
  }

  const handleCopy = () => {
    if (shortUrl) {
      navigator.clipboard.writeText(shortUrl)
        .then(() => {
          setCopySuccess('Copied!')
          setTimeout(() => setCopySuccess(''), 2000) // Clear message after 2 seconds
        })
        .catch(err => {
          console.error('Failed to copy: ', err)
        })
    }
  }

  return (
    <main className="flex min-h-screen flex-col items-center justify-start pt-16 bg-gradient-to-br from-gray-900 via-purple-900 to-violet-900 relative overflow-hidden">      
      <h1 className="text-6xl font-bold text-center bg-clip-text text-transparent bg-gradient-to-r from-blue-400 via-purple-500 to-pink-500 mb-12 animate-bg-gradient">
        Yatiac URL Shortener
      </h1>
      <Card className="w-full max-w-md bg-black bg-opacity-50 border border-purple-500">
        <CardHeader>
          <CardDescription className="text-center text-gray-300">Enter a long URL and we&apos;ll make it shorter for you!</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="url" className="text-gray-300">URL to shorten</Label>
              <div className="relative">
                <Input
                  id="url"
                  type="text"
                  placeholder="https://example.com/very/long/url"
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  className="bg-transparent border-2 border-transparent bg-gradient-to-r from-blue-500 via-purple-500 to-pink-500 rounded-md p-[2px]"
                />
                <div className="absolute inset-[2px] bg-black rounded-md pointer-events-none" />
                <input
                  type="text"
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  className="absolute inset-0 bg-transparent text-white placeholder-gray-400 px-3 py-2 focus:outline-none"
                  placeholder="https://example.com/very/long/url"
                />
              </div>
            </div>
            <Button type="submit" className="w-full bg-gradient-to-r from-blue-500 via-purple-500 to-pink-500 text-white font-bold py-2 px-4 rounded">
              {loading ? (
                
<div role="status">
    <svg aria-hidden="true" className="pr-1 inline w-4 h-4 text-gray-200 animate-spin dark:text-gray-600 fill-blue-600" viewBox="0 0 100 101" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z" fill="currentColor"/>
        <path d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z" fill="currentFill"/>
    </svg>    
</div>

              ) : (
                <Link className="mr-2 h-4 w-4" />
              )}
              {loading ? 'Shortening...' : 'Make it shorter'}
            </Button>
          </form>
          {error && (
            <Alert variant="destructive" className="mt-4 bg-red-900 border-red-500 text-red-200">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
          {shortUrl && (
            <Alert className="mt-4 bg-green-900 border-green-500 text-green-200">
            <Link className="h-4 w-4" />
            <AlertTitle>Success!</AlertTitle>
            <AlertDescription className="mt-2 flex items-center justify-between">
              <span>
                Your shortened URL:{" "}
                <a
                  href={shortUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="font-bold underline text-green-300 hover:text-green-100 transition-colors"
                >
                  {shortUrl}
                </a>
              </span>
              <Button
                onClick={handleCopy}
                className="ml-2 bg-green-700 hover:bg-green-600 text-green-100 font-medium py-1 px-3 rounded transition-colors"
              >
                <Copy className="h-4 w-4 mr-2" />
                {copySuccess ? 'Copied!' : 'Copy'}
              </Button>
            </AlertDescription>
          </Alert>
          )}
        </CardContent>
      </Card>
      <footer className="w-full mt-auto text-center text-gray-400 p-2 bg-gray-800">
        <p>© {new Date().getFullYear()} Yatiac. All rights reserved.</p>
        <p>
          <a href="https://github.com/yatiac" target="_blank" rel="noopener noreferrer" className="text-blue-300">
            <FaGithub className="inline-block h-6 w-6" />
          </a>
        </p>
      </footer>
    </main>
  )
}
