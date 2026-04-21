import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:youtube_rss_ui/service/youtuberss/model/channel.dart';
import 'package:youtube_rss_ui/service/youtuberss/youtuberss.dart';

class UserPage extends StatefulWidget {
  const UserPage({super.key});

  @override
  State<UserPage> createState() => _UserPageState();
}

class _UserPageState extends State<UserPage> {
  final YoutubeRSS _api = YoutubeRSS();
  final int _userId = 1;
  List<Channel> _channels = [];
  bool _isLoading = true;
  int _currentPage = 1;
  final int _pageSize = 12;
  String _url = '';

  @override
  void initState() {
    super.initState();
    _loadUser();
    _loadContent();
  }

  Future<void> _loadUser() async {
    setState(() => _isLoading = true);
    try {
      final channels = await _api.getUserChannels(_userId);
      setState(() {
        _channels = channels;
        _isLoading = false;
      });
    } catch (e) {
      setState(() => _isLoading = false);
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Erro ao carregar usuário: $e')));
      }
    }
  }

  Future<void> _loadContent() async {
    try {
      final pagination = await _api.getUserContentPaginated(
        _userId,
        _currentPage,
        _pageSize,
      );
      setState(() {
        _currentPage = pagination.page;
      });
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erro ao carregar conteúdo: $e')),
        );
      }
    }
  }

  Future<void> _addChannel() async {
    try {
      await _api.addUserChannel(_userId, _url);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Canal adicionado com sucesso!'),
            backgroundColor: Colors.green,
          ),
        );
        await _loadUser();
        await _loadContent();
        setState(() {
          _url = '';
        });
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Erro ao adicionar canal: $e')));
      }
    }
  }

  Future<void> _removeChannel(int channelId) async {
    try {
      await _api.deleteUserChannel(_userId, channelId);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Canal removido com sucesso!'),
            backgroundColor: Colors.green,
          ),
        );
        await _loadUser();
        await _loadContent();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Erro ao remover canal: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    var children = [
      Padding(
        padding: EdgeInsetsGeometry.only(top: 50, left: 16),
        child: Column(
          children: [
            Row(
              children: [
                Container(
                  width: 80,
                  height: 80,
                  decoration: BoxDecoration(
                    color: Colors.yellowAccent,
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: const Center(
                    child: Icon(Icons.person, size: 40, color: Colors.black),
                  ),
                ),
                const SizedBox(width: 16),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Usuário',
                      style: TextStyle(fontSize: 16, color: Colors.grey),
                    ),
                    const SizedBox(height: 4),
                    Text('ID: $_userId', style: const TextStyle(fontSize: 18)),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 16),
            Text(
              'Canais Assinados (${_channels.length})',
              style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
            ),
          ],
        ),
      ),
      _channels.isEmpty
          ? const Padding(
              padding: EdgeInsets.all(16),
              child: Text('Nenhum canal assinado'),
            )
          : SizedBox(
              height: 520,
              child: ListView.builder(
                itemCount: _channels.length,
                itemBuilder: (context, index) {
                  final channel = _channels[index];
                  return _buildChannelCard(channel);
                },
              ),
            ),

      const Divider(height: 32),
      _buildAddChannel(),
    ];
    return Scaffold(
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: children,
              ),
            ),
    );
  }

  Widget _buildAddChannel() {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Adicionar Canal',
            style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 8),
          TextField(
            decoration: InputDecoration(
              hintText: 'Cole o URL do canal do YouTube',
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(8),
              ),
              suffixIcon: IconButton(
                icon: const Icon(Icons.link_rounded),
                onPressed: () {
                  _addChannel();
                },
              ),
            ),
            onChanged: (value) => setState(() {
              _url = value;
            }),
          ),
        ],
      ),
    );
  }

  Widget _buildChannelCard(Channel channel) {
    return Card(
      margin: const EdgeInsets.symmetric(vertical: 4),
      child: ListTile(
        leading: const Icon(
          Icons.play_circle_outline,
          color: Colors.yellowAccent,
        ),
        title: Text(channel.name),
        subtitle: Text(
          channel.rssUrl.substring(52),
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
        ),
        trailing: IconButton(
          icon: const Icon(Icons.delete_outline),
          onPressed: () => _removeChannel(channel.id),
        ),
      ),
    );
  }
}

String formatDate(DateTime date) {
  return DateFormat('dd/MM/yyyy HH:mm').format(date);
}
