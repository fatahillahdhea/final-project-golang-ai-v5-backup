import React, { useState } from "react";
import axios from "axios";

function App() {
  const [file, setFile] = useState(null);
  const [query, setQuery] = useState("");
  const [fileQuery, setFileQuery] = useState("");
  const [response, setResponse] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    setFile(selectedFile);
  };

  const handleFileQueryChange = (e) => {
    setFileQuery(e.target.value);
  };

  const handleQuestionChange = (e) => {
    setQuery(e.target.value);
  };

  const handleUpload = async () => {
    if (!file || !fileQuery) {
      alert("Please select a file and enter a question");
      return;
    }

    setIsLoading(true);
    const formData = new FormData();
    formData.append("file", file);
    formData.append("question", fileQuery);

    try {
      const res = await axios.post('http://localhost:8080/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      setResponse(res.data.answer);
    } catch (error) {
      console.error('Error uploading file:', error);
      setResponse("An error occurred while processing the file.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleChat = async () => {
    if (!query) {
      alert("Please enter a question");
      return;
    }

    setIsLoading(true);
    try {
      const res = await axios.post("http://localhost:8080/chat", { query });
      setResponse(res.data.answer);
    } catch (error) {
      console.error("Error querying chat:", error);
      setResponse("An error occurred while processing the chat.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#f5f5f5] to-[#e0e0e0] flex items-center justify-center p-4">
      <div className="w-full max-w-xl bg-white shadow-2xl rounded-xl overflow-hidden transform transition-all duration-300 hover:scale-[1.02]">
        <div className="bg-gradient-to-r from-[#a7866f] to-[#8d6b4d] p-6">
          <h1 className="text-3xl font-bold text-white text-center flex items-center justify-center gap-3">
            📄 Data Analysis Chatbot
          </h1>
        </div>
        
        <div className="p-6 space-y-4">
          {/* File Upload Section */}
          <div className="space-y-3">
            <div className="flex items-center space-x-3">
              <input 
                type="file" 
                onChange={handleFileChange}
                className="file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:bg-[#a7866f] file:text-white hover:file:bg-[#8d6b4d] transition-all"
              />
              <input 
                type="text" 
                value={fileQuery} 
                onChange={handleFileQueryChange}
                placeholder="Question about the file"
                className="flex-grow p-2 border rounded-md focus:ring-2 focus:ring-[#a7866f] transition-all"
              />
            </div>
            <div className="flex justify-end">
              <button 
                onClick={handleUpload} 
                disabled={isLoading}
                className="bg-[#a7866f] text-white px-4 py-2 rounded-md hover:bg-[#8d6b4d] transition-all flex items-center gap-2 disabled:opacity-50"
              >
                {isLoading ? "Processing..." : "Upload and Analyze"}
              </button>
            </div>
          </div>

          {/* Chat Section */}
          <div className="flex items-center space-x-3">
            <input
              type="text"
              value={query}
              onChange={handleQuestionChange}
              placeholder="Ask a general question..."
              className="flex-grow p-2 border rounded-md focus:ring-2 focus:ring-[#a7866f] transition-all"
            />
            <button 
              onClick={handleChat} 
              disabled={isLoading}
              className="bg-[#a7866f] text-white px-4 py-2 rounded-md hover:bg-[#8d6b4d] transition-all flex items-center gap-2 disabled:opacity-50"
            >
              {isLoading ? "Thinking..." : "Chat"}
            </button>
          </div>

          {/* Response Section */}
          <div 
            className={`mt-6 p-4 border rounded-md bg-gray-50 min-h-[150px] transition-all duration-300 
            ${response ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-4'}
            ${isLoading ? 'animate-pulse' : ''}`}
          >
            <h2 className="text-xl font-semibold mb-3 text-[#5e4b3c]">Response</h2>
            <p className={`${isLoading ? 'text-gray-400' : 'text-gray-700'}`}>
              {isLoading ? "Generating response..." : (response || "Responses will appear here")}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
