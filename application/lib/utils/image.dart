import 'package:camera/camera.dart';
import 'package:image/image.dart' as img;
import 'dart:typed_data';

img.Image yuv420ToRgb(CameraImage image) {
  // 1. Get image dimensions
  int width = image.width;
  int height = image.height;

  // 2. Extract Y, U, and V planes
  final planes = image.planes;

  assert(planes.length == 3);

  // Y plane
  final yBuffer = planes[0].bytes;
  final yRowStride = planes[0].bytesPerRow;
  final yPixelStride = planes[0].bytesPerPixel ?? 1;

  // U (Cb) plane
  final uBuffer = planes[1].bytes;
  final uRowStride = planes[1].bytesPerRow;
  final uPixelStride = planes[1].bytesPerPixel ?? 1;

  // V (Cr) plane
  final vBuffer = planes[2].bytes;
  final vRowStride = planes[2].bytesPerRow;
  final vPixelStride = planes[2].bytesPerPixel ?? 1;

  // 3. Create an image buffer for RGB data (using the image package)
  final rgbImage = img.Image(width: width, height: height);

  // 4. Perform the YUV to RGB conversion and resizing
  for (int y = 0; y < height; y++) {
    for (int x = 0; x < width; x++) {
      final yIndex = y * yRowStride + x * yPixelStride;
      final uIndex =
          (y ~/ 2) * uRowStride + (x ~/ 2) * uPixelStride; // Integer division
      final vIndex =
          (y ~/ 2) * vRowStride + (x ~/ 2) * vPixelStride; // Integer division

      final Y = yBuffer[yIndex];
      final U = uBuffer[uIndex];
      final V = vBuffer[vIndex];

      // YUV to RGB conversion (BT.601 standard is common for camera images)
      int R = (Y + 1.164 * (V - 128)).toInt();
      int G = (Y - 0.391 * (U - 128) - 0.813 * (V - 128)).toInt();
      int B = (Y + 2.018 * (U - 128)).toInt();

      // Clamp RGB values to 0-255 range
      R = R.clamp(0, 255).toInt();
      G = G.clamp(0, 255).toInt();
      B = B.clamp(0, 255).toInt();

      rgbImage.setPixelRgba(x, y, R, G, B, 255);
    }
  }

  return rgbImage;
}

Float32List imageToFloat32List(img.Image image) {
  // Convert to float32 array and normalize
  List<double> imgArray = [];
  // Process channels in RGB order (transpose to channel-first format)
  for (var c = 0; c < 3; c++) {
    for (var y = 0; y < image.height; y++) {
      for (var x = 0; x < image.width; x++) {
        final pixel = image.getPixel(x, y);
        double value;
        switch (c) {
          case 0:
            value = pixel.r.toDouble();
            break;
          case 1:
            value = pixel.g.toDouble();
            break;
          case 2:
            value = pixel.b.toDouble();
            break;
          default:
            value = 0;
        }
        imgArray.add(value / 255.0);
      }
    }
  }
  return Float32List.fromList(imgArray);
}
