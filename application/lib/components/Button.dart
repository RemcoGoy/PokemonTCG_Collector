import 'package:flutter/material.dart';

class TCGButton extends StatelessWidget {

  final Function()? onTap;

  final String buttonText;

  const TCGButton({
    super.key,
    required this.onTap,
    required this.buttonText
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: EdgeInsets.symmetric(horizontal: 25, vertical: 15),
        margin: EdgeInsets.symmetric(horizontal: 25),
        decoration: BoxDecoration(
          color: Colors.grey[300],
          borderRadius: BorderRadius.circular(5),
          border: Border.all(
            color: Colors.deepPurple
          )
        ),
        child: Center(
          child: Text(
            buttonText,
            style: TextStyle(
              color: Colors.deepPurple,
              fontWeight: FontWeight.bold
            ),
          ),
        ),
      ),
    );
  }
}