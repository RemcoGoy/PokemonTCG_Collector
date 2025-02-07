import onnxruntime as ort
import cv2
import numpy as np
from pathlib import Path
import glob

# Initialize ONNX Runtime session
model_path = "../models/best.onnx"
session = ort.InferenceSession(model_path)

# Get model input name
input_name = session.get_inputs()[0].name

# Set up test data directory
test_dir = "../data/annotated/test/images"
test_images = glob.glob(str(Path(test_dir) / "*.jpg"))


def preprocess_image(img_path):
    img = cv2.imread(img_path)
    img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
    img = cv2.resize(img, (640, 640))
    img = img.transpose(2, 0, 1)
    img = img.astype("float32") / 255.0
    return np.expand_dims(img, axis=0), img


# Run inference on test images
for img_path in test_images:
    # Preprocess image
    input_data, orig_img = preprocess_image(img_path)

    # Run inference
    outputs = session.run(None, {input_name: input_data})

    # Process outputs (assuming YOLO v8 output format)
    predictions = outputs[0]  # First output contains detections

    # Load and prepare image for visualization
    img = cv2.imread(img_path)
    img = cv2.resize(img, (640, 640))

    # Draw detections
    for pred in predictions[0]:  # Process first batch
        confidence = pred[4]
        if confidence > 0.25:  # Confidence threshold
            x1, y1, x2, y2 = map(int, pred[:4])

            # Draw bounding box
            cv2.rectangle(img, (x1, y1), (x2, y2), (0, 255, 0), 2)

            # Add confidence label
            label = f"Card: {confidence:.2f}"
            cv2.putText(img, label, (x1, y1-10), cv2.FONT_HERSHEY_SIMPLEX, 0.5, (0, 255, 0), 2)

    # Display image
    cv2.imshow("Detections", img)
    cv2.waitKey(0)  # Wait for key press

cv2.destroyAllWindows()  # Clean up windows when done
