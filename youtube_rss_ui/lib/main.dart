import 'package:flutter/material.dart';
import 'package:youtube_rss_ui/page/video.dart';
import 'package:youtube_rss_ui/page/user.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'Flutter Demo',
      themeMode: ThemeMode.dark,
      darkTheme: ThemeData(
        brightness: Brightness.dark,
        colorScheme: ColorScheme.dark(primary: Colors.yellowAccent),
      ),
      home: const MainPage(),
    );
  }
}

class MainPage extends StatefulWidget {
  const MainPage({super.key});

  @override
  State<MainPage> createState() => _MainPageState();
}

class _MainPageState extends State<MainPage> {
  final List<Widget> _body = [];
  int _currentIndex = 0;

  @override
  void initState() {
    super.initState();
    setState(() {
      _body.addAll(pages());
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      bottomNavigationBar: NavigationBar(
        destinations: navbarbuttons(),
        onDestinationSelected: navbarDestinationSelected,
      ),
      body: _body[_currentIndex],
    );
  }

  void navbarDestinationSelected(int index) {
    setState(() {
      _currentIndex = index;
    });
  }

  List<Widget> navbarbuttons() {
    return const <Widget>[
      NavigationDestination(
        selectedIcon: Badge(
          label: Text('2'),
          child: Icon(Icons.play_arrow_outlined),
        ),
        icon: Icon(Icons.play_arrow),
        label: 'Videos',
      ),
      // NavigationDestination(
      //   icon: Badge(child: Icon(Icons.history)),
      //   label: 'History',
      // ),
      NavigationDestination(
        selectedIcon: Icon(Icons.person),
        icon: Icon(Icons.person),
        label: 'User',
      ),
    ];
  }

  List<Widget> pages() {
    return <Widget>[
      VideoPage(),
      UserPage(),
    ];
  }
}
