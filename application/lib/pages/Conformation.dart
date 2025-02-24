import 'package:application/components/Button.dart';
import 'package:application/main.dart';
import 'package:flutter/material.dart';

class Conformation extends StatelessWidget {
  const Conformation({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Center(
          child: Padding(
            padding: const EdgeInsets.all(35.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Text(
                  'An email has been send to your inbox. Please confirm your email and return to the login page.',
                  style: TextStyle(
                    fontSize: 24,
                  ),
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 50),
                TCGButton(
                  buttonText: 'Return to login',
                  onTap: () => returnButtenPressed(context)
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  void returnButtenPressed(BuildContext context){
    Navigator.push(context, MaterialPageRoute(builder: (context) => MainApp()));
  }
}