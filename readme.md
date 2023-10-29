**ENV Variables Required**

PROJECT_ID --> GCP Project ID to use with Vertex AI

**GCP Auth Required**

The application utilizes Google Default Auth Credentials.

When running in a GCP container such as CloudRun/AppEngine/GKE ensure the service account has permissions to utilize VertexAI.

When running locally ensure your user has permissions to utilize VertexAI and run `gcloud auth application-default login`

**API:**

- `GET /ping` --> healthcheck
- `POST /generate` --> Generate Image using Vertex AI utilizing parameters specified in a JSON. Check out example-request.json and example-response.json as needed.

**Parameter Details**

*prompt*
- The text prompt guides what images the model generates. 
- This field is required.

*sampleImageStyle*
- One of the available predefined styles:
  - photograph
  - digital_art
  - landscape
  - sketch
  - watercolor
  - cyberpunk
  - pop_art

*sampleCount*
- The number of generated images. 
- Accepted integer values: 1-8. 
- Default value: 4.

*negativePrompt*
- A negative prompt to help generate the images.
- For example: 
  - "animals" (removes animals)
  - "blurry" (makes the image clearer)
  - "text" (removes text)
  - "cropped" (removes cropped images)

*seed*
- Any non-negative integer you provide to make output images deterministic. 
- Providing the same seed number always results in the same output images. 
- Accepted integer values: 1 - 2147483647.
