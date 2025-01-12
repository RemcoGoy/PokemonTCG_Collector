import 'package:camera/camera.dart';
import 'package:flutter/material.dart';
import 'package:gal/gal.dart';

class Camera extends StatefulWidget {
  const Camera({super.key});

  @override
  State<Camera> createState() => _CameraState();
}

class _CameraState extends State<Camera> with WidgetsBindingObserver{
  List<CameraDescription> cameras = [];
  CameraController? cameraController;

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);

    if(cameraController == null || cameraController?.value.isInitialized == null){
      return;
    }

    if(state == AppLifecycleState.inactive){
      cameraController!.dispose();
    }else if(state == AppLifecycleState.resumed){
      _setUpCameraController();
    }
  }

  @override
  void initState() {
    super.initState();
    _setUpCameraController();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _buildUI(),
    );
  }

  Widget _buildUI(){
    if(cameraController == null || cameraController?.value.isInitialized == false){
      return const Center(
        child: CircularProgressIndicator()
      );
    }

    return SafeArea(child: SizedBox.expand(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          SizedBox(
            height: MediaQuery.sizeOf(context).height * 0.75,
            width: MediaQuery.sizeOf(context).width,
            child: CameraPreview(cameraController!)
          ),
          IconButton(
            onPressed: () async {
              XFile picture = await cameraController!.takePicture();
              Gal.putImage(picture.path);
            },
            icon: const Icon(
              size: 40,
              Icons.circle_outlined,
              color: Colors.deepPurple
            )
          )
        ],
      ),
    ));

  }

  Future<void> _setUpCameraController() async {
    List<CameraDescription> _cameras = await availableCameras();

    if(_cameras.isNotEmpty){
      setState(() {
        cameras = _cameras;
        cameraController = CameraController(_cameras.first, ResolutionPreset.high);
      });

      cameraController?.initialize().then((_) {
        if(!mounted){
          return;
        }
        setState(() {});
      }).catchError((Object e){
        print(e);
      });
    }
  }
}
