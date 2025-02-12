import 'dart:typed_data';

import 'package:application/utils/image.dart';
import 'package:camera/camera.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:onnxruntime/onnxruntime.dart';
import 'package:image/image.dart' as img;

class DetectionRealtimePage extends StatefulWidget {
  const DetectionRealtimePage({Key? key}) : super(key: key);

  @override
  _DetectionRealtimePageState createState() => _DetectionRealtimePageState();
}

class _DetectionRealtimePageState extends State<DetectionRealtimePage> {
  CameraController? _cameraController;
  bool _isCameraInitialized = false;
  late OrtSession _session;
  List<Map<String, dynamic>>? _detections;
  bool _isProcessing = false;

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
    super.dispose();
  }

  Future<void> _loadModel() async {
    final sessionOptions = OrtSessionOptions();
    const assetFileName = "assets/models/best.onnx";
    final rawAssetFile = await rootBundle.load(assetFileName);
    final bytes = rawAssetFile.buffer.asUint8List();

    // Create ONNX Runtime session
    _session = await OrtSession.fromBuffer(bytes, sessionOptions);
  }

  Float32List _preprocessImage(CameraImage camImage) {
    final format = camImage.format;

    img.Image rgbImage;

    if (format.group == ImageFormatGroup.yuv420) {
      rgbImage = yuv420ToRgb(camImage);
    } else {
      throw UnsupportedError('Image format not supported');
    }

    rgbImage = img.copyResize(rgbImage, width: 640, height: 640);

    return imageToFloat32List(rgbImage);
  }

  Future<void> _processImage(CameraImage camImage) async {
    if (_isProcessing) return;
    _isProcessing = true;

    try {
      final imgArray = _preprocessImage(camImage);

      final inputOrt =
          OrtValueTensor.createTensorWithDataList(imgArray, [1, 3, 640, 640]);
      final inputs = {'images': inputOrt};
      final runOptions = OrtRunOptions();

      try {
        final outputs = await _session.runAsync(runOptions, inputs);

        if (outputs != null) {
          List<List<List<double>>> predictions =
              outputs.first?.value as List<List<List<double>>>;

          setState(() {
            _detections = [];
            for (var elem in predictions) {
              for (var bbox in elem) {
                if (bbox[4] > 0.25) {
                  _detections!.add({
                    'bbox': [bbox[0], bbox[1], bbox[2], bbox[3]],
                    'confidence': bbox[4],
                  });
                }
              }
            }
          });

          // Clean up ONNX resources
          outputs.forEach((element) {
            element?.release();
          });
        }
      } finally {
        // Ensure resources are always released
        inputOrt.release();
        runOptions.release();
      }
    } catch (e) {
      print('Error processing image: $e');
    } finally {
      _isProcessing = false;
    }
  }

  Future<void> _initializeCamera() async {
    final cameras = await availableCameras();
    if (cameras.isEmpty) return;

    _cameraController = CameraController(
      cameras.first,
      ResolutionPreset.medium,
      enableAudio: false,
    );

    await _cameraController!.initialize();

    _cameraController!.startImageStream((image) {
      if (!mounted) return;
      _processImage(image);
    });

    if (mounted) {
      setState(() {
        _isCameraInitialized = true;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Stack(
        children: [
          _isCameraInitialized
              ? CameraPreview(_cameraController!)
              : const Center(
                  child: CircularProgressIndicator(),
                ),
          Positioned(
            top: 40,
            left: 20,
            child: Container(
              padding: const EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: Colors.black54,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Text(
                'Detections: ${_detections?.length ?? 0}',
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
