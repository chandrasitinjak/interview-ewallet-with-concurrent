1. Bagaimana Anda memastikan bahwa operasi kredit dan debit bersifat atomic?
->Semua query (GetUserForUpdate, UpdateBalanceTx, CreateTransactionTx) dijalankan di dalam satu transaksi DB (tx).
Kalau salah satu step gagal, defer tx.Rollback() otomatis membatalkan semua perubahan. Kalau semua sukses → tx.Commit() baru dijalankan.
2. Jelaskan potensi race condition yang mungkin terjadi dalam sistem ini dan bagaimana Anda mengatasinya.
->Dua goroutine bisa mengakses user balance bersamaan. Kalau tidak dikontrol, hasil akhir bisa salah.
solusi : pakai sync.Mutex (s.mu.Lock()) di service dan pakai SELECT ... FOR UPDATE di GetUserForUpdate, artinya row user dikunci di level DB sampai transaksi selesai. Ini mencegah dua transaksi DB update saldo user yang sama di waktu bersamaan.
3.Bagaimana Anda akan mengimplementasikan rollback mekanisme jika terjadi kegagalan di tengah proses transaksi?
-> pakai defer tx.Rollback(), sehingga kalau ada error sebelum Commit(), semua perubahan otomatis dibatalkan.