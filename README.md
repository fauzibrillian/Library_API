# Library_API
MVP
- user
daftar
login
detail user ( admin + user yang terkait )
search + all user ( halaman admin ) 
delete user ( admin + user yang terkait )
update user ( admin + user yang terkait )

- book 
hanya admin dapat menambah + edit + menghapus book
semua role dapat melihat book ( pagination + searching )
manajemen rak book ( terdapat kode rak buku )

- peminjaman
user harus setor jaminan saat meminjam book ( ktp , ktm , uang )
user dapat melihat list buku yang dipinjam -> ada detail data peminjaman + penalti
menghasilkan nota pengembalian berisi rak jaminan, book apa saja yang dipinjam dan tanggal pengembalian
admin dapat monitoring semua pinjaman

- penalti
terdapat pengingat sms ke nomer user
terdapat kelipatan 5000 per hari jika telat mengembalikan buku
pembayaran bisa pakek midtrans