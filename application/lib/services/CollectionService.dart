import 'package:application/models/Collection.dart';
import 'package:application/testData/Collections.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

import 'package:application/models/pokemon.dart';

Future<List<Pokemon>> fetchCardsForCollection() async { 
  final uri = Uri.parse('https://api.pokemontcg.io/v2/cards?pageSize=50');
  final response = await http.get(uri);

  if(response.statusCode == 200){
    final data = json.decode(cardsCollection1)['data'];
    return List.generate(data.length, (index){
      return Pokemon.fromJson(data[index]);
    });
  }else{
    throw Exception('Something whent wrong while fetching the cards');
  }
}

Future<List<Collection>> fetchCollections() async {
  final uri = Uri.parse('https://api.pokemontcg.io/v2/cards?pageSize=1');
  final response = await http.get(uri);

  if(response.statusCode == 200){
    final data = json.decode(collections)['data'];
    return List.generate(data.length, (index){
      return Collection.fromJson(data[index]);
    });
  }else{
    throw Exception('Something whent wrong while fetching the collection');
  }
}

