import 'dart:convert';

import 'package:application/models/pokemon.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class Home extends StatefulWidget {
  const Home({super.key});

  @override
  State<Home> createState() => _ListState();
}

class _ListState extends State<Home>{

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<List<Pokemon>>(
          future: fetch(),
          builder: (context, snapshot) {
            if(snapshot.connectionState == ConnectionState.waiting){
              return const CircularProgressIndicator();
            }else if(snapshot.connectionState == ConnectionState.none){
              return Container();
            }else{
              if(snapshot.hasData){
                return buildDataWidget(context, snapshot);
              }else if (snapshot.hasError){
                return Text('${snapshot.error}');
              }else{
                return Container();
              }
            }
          }
        );
  }


  Widget buildDataWidget(context, snapshot) => GridView.count(
  crossAxisCount: 2,
  children: List.generate(snapshot.data.length, (index) {
    return Column(
      children: [
        Image.network(
          snapshot.data[index].imageUrl,
          height: 170,
        ),
        Text(snapshot.data[index].name)
      ]
    );
  }),
);

  Future<List<Pokemon>> fetch() async {
    final uri = Uri.parse('https://api.pokemontcg.io/v2/cards?pageSize=50');
    final response = await http.get(uri);

    if(response.statusCode == 200){
      final data = json.decode(response.body)['data'];
      return List.generate(data.length, (index){
        return Pokemon.fromJson(data[index]);
      });
    }else{
      throw Exception('Fuck!!');
    }
  }
}