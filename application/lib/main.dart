import 'package:application/components/Appbar.dart';
import 'package:application/pages/Home.dart';
import 'package:application/pages/Profile.dart';
import 'package:application/pages/camera.dart';
import 'package:application/pages/login.dart';
import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

void main() async {
  
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize Hive
  await Hive.initFlutter();
  await Hive.openBox('userStorage');

  // Initialize .env
  await dotenv.load(fileName: '.env');

  runApp(const MainApp());
}

class MainApp extends StatefulWidget {
  const MainApp({super.key});

  @override
  State createState() => _HomeState();
}

class _HomeState extends State{
  final userStorage = Hive.box('userStorage');

  int _selectedIndex = 0;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      theme: ThemeData(fontFamily: 'red hat text'),
      home: loadHome(),
    );
  }

  Widget loadHome() {
    if(!this.userStorage.containsKey('email') || !this.userStorage.containsKey('AuthToken')){
      return Login();
    }
    
    return Scaffold(
      appBar: AppBar(title: PokeAppBar()),
      bottomNavigationBar: BottomNavigationBar(
        unselectedItemColor: Colors.white70,
        selectedItemColor: Colors.deepPurpleAccent,
        selectedLabelStyle: TextStyle(fontWeight: FontWeight.bold),
        backgroundColor: Colors.black87,
        elevation: 0,
        currentIndex: _selectedIndex,
        onTap: _onItemTapped,
        items: [
          BottomNavigationBarItem(
            icon: Icon(
              Icons.home,
            ),
            label: 'Home'
          ),
          BottomNavigationBarItem(
            icon: Icon(
              Icons.camera,
            ),
            label: 'Scan'
          ),
          BottomNavigationBarItem(
            icon: Icon(
              Icons.person,
            ),
            label: 'Profile'
          ),
        ]
      ),
      body: Center(
        child: _pages.elementAt(_selectedIndex),
      ),
    );
  }

  static const List<Widget> _pages = <Widget>[
    Home(),
    Camera(),
    Profile()
  ];

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }
}