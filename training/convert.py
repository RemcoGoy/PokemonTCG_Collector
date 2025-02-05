from ultralytics import YOLO

# Load the trained YOLO model
model = YOLO('runs/train/pokemon_detector3/weights/best.pt')  # Path to your trained model

# Export the model to ONNX format
model.export(format='onnx',
            imgsz=640,     # Image size used during training
            simplify=True,  # Simplify the model where possible
            opset=12)      # ONNX opset version

print("Model converted to ONNX format. Saved as 'runs/train/pokemon_detector/weights/best.onnx'")
