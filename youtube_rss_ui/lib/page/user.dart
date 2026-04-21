import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:intl/intl.dart';
import 'package:youtube_rss_ui/service/youtuberss/model/channel.dart';
import 'package:youtube_rss_ui/service/youtuberss/model/content.dart';
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
  List<Content> _contents = [];
  bool _isLoading = true;
  final bool _isLoadingMore = false;
  int _currentPage = 1;
  int _totalPages = 1;
  final int _pageSize = 12;

  @override
  void initState() {
    super.initState();
    _loadUser();
    _loadContent();
  }

  Future<void> _loadUser() async {
    setState(() => _isLoading = true);
    try {
      final channels = await _api.GetUserChannels(_userId);
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
        _contents = pagination.items;
        _currentPage = pagination.page;
        _totalPages = pagination.totalPages;
      });
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erro ao carregar conteúdo: $e')),
        );
      }
    }
  }

  Future<void> _addChannel(String channelUrl) async {
    try {
      await _api.addUserChannel(_userId, channelUrl);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Canal adicionado com sucesso!'),
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
        ).showSnackBar(SnackBar(content: Text('Erro ao adicionar canal: $e')));
      }
    }
  }

  Future<void> _removeChannel(int channelId) async {
    try {
      await _api.DeleteUserChannel(_userId, channelId);
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

  Future<void> _updateContent() async {
    setState(() => _isLoading = true);

    await Future.delayed(Duration(seconds: 1));

    try {
      final pagination = await _api.getUserContentPaginated(
        _userId,
        1,
        _pageSize,
      );
      setState(() {
        _contents = pagination.items;
        _currentPage = pagination.page;
        _totalPages = pagination.totalPages;
        _isLoading = false;
      });
    } catch (e) {
      setState(() => _isLoading = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erro ao atualizar conteúdo: $e')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Meu Perfil'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _updateContent,
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Profile Section
                  Padding(
                    padding: const EdgeInsets.all(16),
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
                                child: Icon(Icons.person, size: 40),
                              ),
                            ),
                            const SizedBox(width: 16),
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'Usuário',
                                  style: TextStyle(
                                    fontSize: 16,
                                    color: Colors.grey,
                                  ),
                                ),
                                const SizedBox(height: 4),
                                Text(
                                  'ID: $_userId',
                                  style: const TextStyle(fontSize: 18),
                                ),
                              ],
                            ),
                          ],
                        ),
                        const SizedBox(height: 16),
                        Text(
                          'Canais Assinados (${_channels.length})',
                          style: const TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),
                  ),
                  _channels.isEmpty
                      ? const Padding(
                          padding: EdgeInsets.all(16),
                          child: Text('Nenhum canal assinado'),
                        )
                      : ListView.builder(
                          shrinkWrap: true,
                          physics: const NeverScrollableScrollPhysics(),
                          itemCount: _channels.length,
                          itemBuilder: (context, index) {
                            final channel = _channels[index];
                            return Card(
                              margin: const EdgeInsets.symmetric(
                                horizontal: 16,
                                vertical: 4,
                              ),
                              child: ListTile(
                                leading: Icon(
                                  Icons.play_circle_outline,
                                  color: Colors.yellowAccent,
                                ),
                                title: Text(channel.name),
                                subtitle: Text(
                                  channel.rssUrl,
                                  maxLines: 1,
                                  overflow: TextOverflow.ellipsis,
                                ),
                                trailing: IconButton(
                                  icon: const Icon(Icons.delete_outline),
                                  onPressed: () =>
                                      _removeChannel(channel.id ?? 0),
                                ),
                              ),
                            );
                          },
                        ),

                  const Divider(height: 32),

                  // Add Channel Section
                  Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(
                          'Adicionar Canal',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 8),
                        TextField(
                          decoration: InputDecoration(
                            hintText: 'Cole o URL do canal do YouTube',
                            prefixIcon: const Icon(Icons.link),
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(8),
                            ),
                          ),
                          onChanged: (value) => _addChannel(value),
                        ),
                      ],
                    ),
                  ),

                  const Divider(height: 32),

                  // Content Section
                  Padding(
                    padding: const EdgeInsets.all(16),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          'Conteúdo Recente',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        Text(
                          '$_totalPages páginas',
                          style: const TextStyle(color: Colors.grey),
                        ),
                      ],
                    ),
                  ),

                  _contents.isEmpty
                      ? const Padding(
                          padding: EdgeInsets.all(16),
                          child: Text('Nenhum conteúdo encontrado'),
                        )
                      : ListView.builder(
                          shrinkWrap: true,
                          physics: const NeverScrollableScrollPhysics(),
                          padding: const EdgeInsets.symmetric(vertical: 8),
                          itemCount:
                              _contents.length + (_isLoadingMore ? 1 : 0),
                          itemBuilder: (context, index) {
                            if (index == _contents.length) {
                              return const Center(
                                child: Padding(
                                  padding: EdgeInsets.all(16),
                                  child: CircularProgressIndicator(),
                                ),
                              );
                            }
                            final content = _contents[index];
                            return _buildContentCard(content);
                          },
                        ),

                  _isLoadingMore
                      ? const Padding(
                          padding: EdgeInsets.all(16),
                          child: Center(child: CircularProgressIndicator()),
                        )
                      : const SizedBox.shrink(),
                ],
              ),
            ),
    );
  }

  Widget _buildContentCard(Content content) {
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
      child: InkWell(
        onTap: () async {
          await openUrl(content.url);
        },
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            ClipRRect(
              borderRadius: const BorderRadius.vertical(
                top: Radius.circular(8),
              ),
              child: Image.network(
                content.image.replaceFirst("hqdefault", "maxresdefault"),
                width: double.infinity,
                height: 180,
                fit: BoxFit.cover,
                errorBuilder: (context, error, stackTrace) => Container(
                  color: Colors.grey[800],
                  child: const Icon(Icons.video_library, color: Colors.grey),
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    content.title,
                    style: const TextStyle(fontWeight: FontWeight.w600),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 4),
                  Text(
                    formatDate(content.date),
                    style: const TextStyle(color: Colors.grey),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

String formatDate(DateTime date) {
  return DateFormat('dd/MM/yyyy HH:mm').format(date);
}

Future<void> openUrl(String url) async {
  final uri = Uri.parse(url.replaceFirst('www.youtube.com', 'm.youtube.com'));
  if (await canLaunchUrl(uri)) {
    await launchUrl(uri, mode: LaunchMode.externalApplication);
  }
}
