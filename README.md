# 🌈 bifrost

A tool designed to facilitate magnet link downloading directly onto your [Plex](https://www.plex.tv/) media server.

## 🚀 Introduction

Bifrost is a Go web server specially crafted for remotely managing and downloading files straight into your media server - making your media management experience seamless. Whether it’s movies, TV shows, or any other media content, easily specify the download details, and the server takes care of the rest, placing your files exactly where you want them in your Plex server.

## 🛠 Installation

### Prerequisites

- **Go:** Make sure you have Go installed on your server. If not, you can install it by following the instructions [here](https://golang.org/doc/install).

### Steps

1. **Clone the Repository:**
   Navigate to the desired location on your server and run:

   ```sh
   git clone git@github.com:jacogrande/bifrost.git
   cd bifrost
   ```

2. **Run the Server:**
   ```sh
   ./bifrost
   ```

## 🚩 Options

### --init flag

Bifrost supports multiple download destinations, allowing you to have separate directories for your movies, tv shows, and more.
To make the setup process smooth, an initialization option `--init` is available which prompts user inputs to configure download folders, ensuring a quick and easy setup.

#### How it Works

- **Get Folder Name and Path:**
  You will be prompted to enter a name and a path for the folder where the downloaded content should be stored.
- **Permission Check:**
  The provided folder path will be checked for read and write permissions to ensure smooth file storage.
- **Save Folder Information:**
  Upon successful verification, the folder’s name and path will be stored, and used for subsequent downloads.

This function loopingly prompts you to enter folder details and saves them upon validating the provided path's permissions.

## 📥 /download Endpoint

### Overview

The `/download` endpoint is your main interface for initiating downloads to your media server.

### API Details

**Endpoint:** `/download`

**Method:** `POST`

**Request Payload:**

- **magnet:** The magnet link of the torrent.
- **folder:** The name of the folder (configured during --init) where the downloaded file should be stored.
- **name:** A custom name for the downloaded file.
- **posterUrl:** A URL for the poster image.

#### Example Request

```json
{
  "magnet": "magnet:?xt=urn:btih:example",
  "folder": "MyMovies",
  "name": "Example Movie",
  "posterUrl": "https://example.com/poster.jpg"
}
```

This endpoint enables you to specify exactly where and how the downloaded content should be stored and managed within your media server, ensuring an organized and straightforward user experience.
