import express from "express";
import { createProxyMiddleware } from "http-proxy-middleware";
import path from "path";
import dotenv from "dotenv";

dotenv.config();
const app = express();
const PORT = process.env.PORT || 5000;

// Reverse Proxy to another backend
app.use(
  "/api",
  createProxyMiddleware({
    target: process.env.BACKEND_URL, // The actual backend service
    changeOrigin: true,
  })
);

// Serve static files from React frontend
const __dirname = path.resolve();
app.use(express.static(path.join(__dirname, "../frontend/dist")));

app.get("*", (req, res) => {
  res.sendFile(path.join(__dirname, "../frontend/dist/index.html"));
});

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});

