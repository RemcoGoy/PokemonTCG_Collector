from ultralytics import YOLO
import yaml

# Define paths
DATASET_PATH = '../data/annotated'  # Path to dataset directory
YAML_PATH = '../data/annotated/data.yaml'   # Path to data.yaml file

# Create data.yaml configuration
data = {
    'path': DATASET_PATH,
    'train': 'train/images',  # Training images folder
    'val': 'valid/images',    # Validation images folder
    'test': 'test/images',    # Test images folder
    'names': ['card']      # Class names
}

# Write data.yaml file
with open(YAML_PATH, 'w') as f:
    yaml.dump(data, f)

# Initialize YOLO model
model = YOLO('yolov8n.pt')  # Load pretrained YOLOv8 nano model

# Train the model
results = model.train(
    data=YAML_PATH,
    epochs=100,              # Number of training epochs
    imgsz=640,              # Image size
    batch=16,               # Batch size
    workers=4,              # Number of worker threads
    device='cpu',             # GPU device (use 'cpu' if no GPU available)
    project='runs/train',   # Project name
    name='pokemon_detector' # Experiment name
)

# Validate the model
results = model.val()

print("Training completed. Model saved in runs/train/pokemon_detector")
