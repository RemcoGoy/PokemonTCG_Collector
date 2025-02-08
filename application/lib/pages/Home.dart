import 'package:application/models/Collection.dart';
import 'package:application/models/pokemon.dart';
import 'package:application/services/CollectionService.dart';
import 'package:flutter/material.dart';

class Home extends StatefulWidget {
  const Home({super.key});

  @override
  State<Home> createState() => _ListState();
}

class _ListState extends State<Home>{
  String? dropdownValue = 'One';

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Expanded(
          flex: 1,
          child: FutureBuilder<List<Pokemon>>(
            future: fetchCardsForCollection(),
            builder: (context, snapshot) {
              if(snapshot.connectionState == ConnectionState.waiting){
                return Center(
                  child: const CircularProgressIndicator()
                );
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
          ),
        ),
        Expanded(
          flex: 0,
          child: FutureBuilder<List<Collection>>(
            future: fetchCollections(),
            builder: (context, snapshot) {
              if(snapshot.connectionState == ConnectionState.waiting){
                return Container();
              }else if(snapshot.connectionState == ConnectionState.none){
                return Container();
              }else{
                if(snapshot.hasData){
                  return buildDropdown(context, snapshot);
                }else if (snapshot.hasError){
                  return Text('${snapshot.error}');
                }else{
                  return Container();
                }
              }
            }
          ),
        )
      ]
    );
  }

  Widget buildDropdown(context, snapshot) {
    return Container(
      padding: const EdgeInsets.only(right: 10),
      alignment: Alignment.bottomRight,
      child: DropdownButton<String>(
        value: dropdownValue,
        icon: const Icon(Icons.arrow_downward),
        onChanged: (String? newValue){
          setState(() {
            dropdownValue = newValue;
          });
        },
        items: [],
      ),
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
}