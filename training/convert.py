from ultralytics import YOLO

# Load the trained YOLO model
model = YOLO('../models/best.pt')

# Export the model to ONNX format
model.export(
    format='onnx',
    nms=True,
    opset=18
)

print("Model converted to ONNX format. Saved as 'runs/train/pokemon_detector/weights/best.onnx'")
