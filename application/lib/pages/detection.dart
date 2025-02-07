import 'dart:io';
import 'dart:typed_data';
import 'package:flutter/material.dart';
import 'package:camera/camera.dart';
import 'package:flutter/services.dart';
import 'package:image/image.dart' as img;
import 'package:onnxruntime/onnxruntime.dart';

class DetectionPage extends StatefulWidget {
  const DetectionPage({Key? key}) : super(key: key);

  @override
  _DetectionPageState createState() => _DetectionPageState();
}

class _DetectionPageState extends State<DetectionPage> {
  File? _image;
  List<Map<String, dynamic>>? _detections;
  late OrtSession _session;
  CameraController? _cameraController;
  bool _isCameraInitialized = false;

  @override
  void initState() {
    super.initState();
    OrtEnv.instance.init();
    _loadModel();
    _initializeCamera();
  }

  @override
  void dispose() {
    _cameraController?.dispose();
    OrtEnv.instance.release();
    super.dispose();
  }

  Future<void> _initializeCamera() async {
    final cameras = await availableCameras();
    if (cameras.isEmpty) return;

    _cameraController = CameraController(
      cameras.first,
      ResolutionPreset.high,
      enableAudio: false,
    );

    await _cameraController!.initialize();
    setState(() {
      _isCameraInitialized = true;
    });
  }

  Future<void> _loadModel() async {
    final sessionOptions = OrtSessionOptions();
    const assetFileName = "assets/models/best.onnx";
    final rawAssetFile = await rootBundle.load(assetFileName);
    final bytes = rawAssetFile.buffer.asUint8List();

    // Create ONNX Runtime session
    _session = await OrtSession.fromBuffer(bytes, sessionOptions);
  }

  Future<void> _takePhoto() async {
    if (!_isCameraInitialized || _cameraController == null) return;

    final XFile photo = await _cameraController!.takePicture();
    setState(() {
      _image = File(photo.path);
      _processImage();
    });
  }

  Future<void> _processImage() async {
    if (_image == null) return;

    // Load and preprocess image
    final imageBytes = await _image!.readAsBytes();
    final decodedImage = img.decodeImage(imageBytes);
    if (decodedImage == null) return;

    // Resize to 640x640
    final resizedImage = img.copyResize(decodedImage, width: 640, height: 640);

    // Convert to float32 array and normalize
    List<double> imgArray = [];
    // Process channels in RGB order (transpose to channel-first format)
    for (var c = 0; c < 3; c++) {
      for (var y = 0; y < 640; y++) {
        for (var x = 0; x < 640; x++) {
          final pixel = resizedImage.getPixel(x, y);
          double value;
          switch (c) {
            case 0: value = pixel.r.toDouble(); break;
            case 1: value = pixel.g.toDouble(); break;
            case 2: value = pixel.b.toDouble(); break;
            default: value = 0;
          }
          imgArray.add(value / 255.0);
        }
      }
    }
    final inputData = Float32List.fromList(imgArray);

    // Run inference
    final inputOrt = OrtValueTensor.createTensorWithDataList(inputData, [1, 3, 640, 640]);
    final inputs = {'images': inputOrt};
    final runOptions = OrtRunOptions();
    final outputs = await _session.runAsync(runOptions, inputs);
    List<List<List<double>>> predictions = outputs?.first?.value as List<List<List<double>>>;

    setState(() {
      _detections = [];
      predictions.forEach((elem) {
        elem.forEach((bbox) {
          if (bbox[4] > 0.25) {
            _detections!.add({
              'bbox': [
                bbox[0], // x1
                bbox[1], // y1
                bbox[2], // x2
                bbox[3], // y2
              ],
              'confidence': bbox[4],
            });
          }
        });
      });
    });

    inputOrt.release();
    runOptions.release();
    outputs?.forEach((element) {
      element?.release();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Card Detection'),
      ),
      body: Column(
        children: [
          if (!_isCameraInitialized) ...[
            const Expanded(
              child: Center(
                child: CircularProgressIndicator(),
              ),
            ),
          ] else if (_image == null) ...[
            Expanded(
              child: CameraPreview(_cameraController!),
            ),
          ] else ...[
            Expanded(
              child: CustomPaint(
                child: Image.file(_image!),
                foregroundPainter: BoundingBoxPainter(
                  detections: _detections ?? [],
                ),
              ),
            ),
          ],
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _image == null ? _takePhoto : () {
          setState(() {
            _image = null;
            _detections = null;
          });
        },
        child: Icon(_image == null ? Icons.camera_alt : Icons.refresh),
      ),
    );
  }
}

class BoundingBoxPainter extends CustomPainter {
  final List<Map<String, dynamic>> detections;

  BoundingBoxPainter({required this.detections});

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = Colors.green
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2.0;

    for (final detection in detections) {
      final bbox = detection['bbox'] as List<double>;
      final rect = Rect.fromLTRB(
        bbox[0],
        bbox[1],
        bbox[2],
        bbox[3],
      );
      canvas.drawRect(rect, paint);

      // Draw confidence text
      final textPainter = TextPainter(
        text: TextSpan(
          text: 'Card: ${(detection['confidence'] * 100).toStringAsFixed(1)}%',
          style: const TextStyle(
            color: Colors.green,
            fontSize: 16,
            backgroundColor: Colors.black54,
          ),
        ),
        textDirection: TextDirection.ltr,
      );
      textPainter.layout();
      textPainter.paint(canvas, Offset(rect.left, rect.top - 20));
    }
  }

  @override
  bool shouldRepaint(CustomPainter oldDelegate) => true;
}
