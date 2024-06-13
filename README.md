# Pinjol API

#### Final project Enigma Camp

Built using [go](https://go.dev/) v1.22.0

## API Spec

### Users API
## Pebri
#### Login

Request :

- Method : `POST`
- Endpoint : `/users/login`
- Body :

```json
{
  "email": "string",
  "password": "string"
}
```

Response :

- Status : 200 Created
- Body :

```json
{
  "responseCode": "string",
  "data": {
    "token": "string"
  }
}
```

#### Create User

Request :

- Method : `POST`
- Endpoint : `/users/create`
- Body :

```json
{
  "fullName": "string",
  "email": "string",
  "password": "string"
}
```

Response :

- Status : 200 Created
- Body :

```json
{
  "responseCode": "string",
  "message": "message"
}
```

## ?
#### Isi Data Diri User Debitur

Request :

- Method : `POST`
- Endpoint : `/users/debitur/:id/form`
- Body :

```json
{
  "detail": {
     "nik": "string",
     "phoneNumber": "string",
     "address": "string",
     "city": "string",
     "fotoKTP": "string",
     "fotoSelfie": "string",
  },
  "jobs": {
     "jobName": "string",
     "gaji": number,
     "officeName": "string",
     "officeContact": "string",
     "officeAddress": "string",
  },
  "emergency": {
     "contactName": "string",
     "contactNumber": "string"
  }
}
```

Response :

- Status : 200 OK
- Body :

```json
{
  "responseCode": "string",
  "message": "message"
}
```

#### Isi Data Diri User Debt Collector

Request :

- Method : `POST`
- Endpoint : `/users/debt-collector/:id/form`
- Body :

```json
{
  "nik": "string",
  "phoneNumber": "string",
  "address": "string",
  "city": "string",
  "fotoKTP": "string",
  "fotoSelfie": "string"
}
```

Response :

- Status : 200 OK
- Body :

```json
{
  "responseCode": "string",
  "message": "message"
}
```

#### Get All Data By Roles (debitur/debt-collector)

Request :

- Method : `GET`
- Endpoint : `/users/:roles`
- Body :
- QueryParam: `?page=1&size=5&?status=verified`
- QueryParam: `?page=1&size=5&?status=unverified`

Response :

- Status : 200
- Body :

```json
{
  "responseCode": 200,
  "data": [
    {
      "userId": "int",
      "NIK": "string",
      "fullname": "string",
      "phoneNumber": "string",
      "address": "string",
      "fotoKtp": "string",
      "fotoSelfie": "string"
    }
  ],
  "paging": {
    "page": 1,
    "totalData": 10
  }
}
```

#### Get All Data User (query:verified,unverified)

Request :

- Method : `GET`
- Endpoint : `/users/`

Response :

- Status : 200 OK
- Body :

```json
{
  "responseCode": "string",
  "data": [
    {
      "id": "string",
      "roles": "string",
      "nik": "string",
      "fullName": "string",
      "phoneNumber": "string",
      "address": "string",
      "city": "string",
      "fotoKTP": "string",
      "fotoSelfie": "string"
    }
  ]
}
```

#### Get Debitur By Id

Request :

- Method : `GET`
- Endpoint : `/users/debitur/:id`

Response :

- Status : 200 OK
- Body :

```json
{
  "responseCode": "string",
  "data":
      {
         "nik": "string",
         "phoneNumber": "string",
         "address": "string",
         "city": "string",
         "fotoKTP": "string",
         "fotoSelfie": "string",
      },
      "jobs": {
         "jobName": "string",
         "gaji": number,
         "officeName": "string",
         "officeContact": "string",
         "officeAddress": "string",
      },
      "emergency": {
         "contactName": "string",
         "contactNumber": "string"
      }
}
```

#### Get Debt Collector By Id

Request :

- Method : `GET`
- Endpoint : `/users/debt-collector/:id`

Response :

- Status : 200 OK
- Body :

```json
{
  "responseCode": "string",
  "data": {
    "nik": "string",
    "phoneNumber": "string",
    "address": "string",
    "city": "string",
    "fotoKTP": "string",
    "fotoSelfie": "string"
  }
}
```

#### Update Data Debitur

- Method : `PUT`
- Endpoint : `/users/debitur/:id/form`
- Body :

```json
{
  "email": "string",
  "password": "string",
  "fullName": "string",
  "phoneNumber": "string",
  "address": "string",
  "city": "string"
}
```

Response :

- Status : 200
- Body :

```json
{
  "responseCode": "string",
  "message": "message"
}
```

#### Update Data Debt Collector

- Method : `PUT`
- Endpoint : `/users/debt-collector/:id/form`
- Body :

```json
{
  "email": "string",
  "password": "string",
  "fullName": "string",
  "phoneNumber": "string",
  "address": "string",
  "city": "string"
}
```

Response :

- Status : 200
- Body :

```json
{
  "responseCode": "string",
  "message": "message"
}
```

#### Isi Nomor Rekening User (Debitur/Debt Collector)

Request :

- Method : Post
- Endpoint : `/users/:id/rekening`
- Body :

```json
{
  "accountNumber": "string",
  "bankName": "string"
}
```

Response :

- Status : 200 OK
- Body :

```json
{
  "responseCode": "string",
  "message": "message"
}
```

## Billy
### Admin

#### Verifikasi Akun by Admin
Request :

- Method : `PATCH`
- Endpoint : `/users/{id}/status`
```json
{
  "status":"verified"
}
```
Response :

- Status : 201 Ok
- Body :

```json
{
  "responseCode":"string",
  "data": {
    "user"
  }
}
```

#### Verifikasi Pinjaman by Admin + Mengirim uang pinjaman
Request :

- Method : `PATCH`
- Endpoint : `/loans/{id}/status`
```json
{
  "status":"verified"
}
```
Response :

- Status : 201 Ok
- Body :

```json
{
  "responseCode":"string",
  "data": {
    "user"
  }
}
```

#### Verifikasi Tarik Uang Debt Collector by Admin + Mengirim uang gaji
- Method : `PATCH`
- Endpoint : `/debt-collector/tugas/{id}/status`
```json
{
  "status":"done"
}
```
Response :

- Status : 201 Ok
- Body :

```json
{
  "responseCode":"string",
  "data": {
    "user"
  }
}
```

## Doni
### Debitur

#### Pengajuan Pinjaman

Request :

- Method : `Post`
- Endpoint : `/users/debitur/create/pinjaman`
- Body :

```json
{
  "jumlahPinjaman": float,
  "tenor": int,
  "description": "string"
}
```

Response :

- Status : 200 OK
- Body :

```json
{
  "responseCode": "string",
  "message": "message"
}
```

#### Get Riwayat Pinjaman

Request :

- Method : `GET`
- Endpoint : `/users/debitur/pinjaman/:id`

Response:

- Status : 200 OK
- Body:

```json
{
  "responseCode": "string",
  "data":
  [
    {
      "userId": int,
      "jumlahPinjaman": float,
      "tenor": int,
      "bungaPerBulan": float,
      "description": "string",
      "statusPengajuan": "string enum",
      "startDate": "string",
      "endDate": "string"
    },
    {
      "userId": int,
      "jumlahPinjaman": float,
      "tenor": int,
      "bungaPerBulan": float,
      "description": "string",
      "statusPengajuan": "string enum",
      "startDate": "string",
      "endDate": "string"
    }
  ]
}
```

#### Get Daftar Cicilan (query:dibayar/belum)

Request:

- Method : `GET`
- Endpoint : `/users/debitur/cicilan/:id`
- Query Params:
  ` ?page=1&size=5`
  `?page=1&size=5&status=dibayar/belum`
- Body :

```json
"responseCode": "string",
"data":{
  "pinjamanId": id,
  "tanggalJatuhTempo": "string",
  "tanggalSelesaiBayar": "string",
  "jumlahBayar": float,
  "status": "string"
}
```

#### Pembayaran Cicilan

Request:

- Method : `POST`
- Endpoint : `/users/debitur/cicilan/pay`

Request:

- Body:

```json
{
  "pinjamaId": int,
  "jumlahBayar": float
}
```

Response:

- Body:

```json
"responseCode": "string",
"data":{
  "pinjamaId": int,
  "tanggalJatuhTempo": "string",
  "tanggalBayarSelesai": "string",
  "jumlahBayar": float,
  "status": "string"
}
```

## Oho
### Debt Collector

#### Get Debitur Nunggak (Mengambil userId(debitur) yang cicilan nunggak > 2 bulan & Belum di claim DC lain)

Request :

- Method: `GET`
- Endpoint: `/debt-collector/debitur-nunggak`

Response :

- Status : 200 OK
- Body:

```json
{
  "status": "success",
  "debiturNunggak": [
    {
      "userId": "string",
      "fullNama": "string",
      "address": "string",
      "cicilanTertunggak": number
    },
    {
      "user_id": "string",
      "nama": "string",
      "cicilanTertunggak": number
    }
  ]
}
```

#### Claim Tugas Debt Collector (Mengambil userId(debitur) yang cicilan nunggak > 2 bulan & Belum di claim DC lain)

Request :

- Method: `POST`
- Endpoint: `/debt-collector/claim-tugas`

- Body :

```json
{
  "user_id": "string"
}
```

Response :

- Status : 200 OK
- Body:

```json
{
  "status": "success",
  "message": "Tugas untuk debitur dengan user_id: 123 berhasil di-claim."
}
```

#### Melihat Gaji

Request :

- Method: `GET`
- Endpoint: `/debt-collector/melihat-gaji/:id`

Response :

- Status : 200 OK
- Body :

```json
{
  "user_id": "string",
  "fullName": "string",
  "totalSalary": "string",
  "historyTugas": [
    {
      "userId": "string",
      "fullNama": "string",
      "address": "string",
      "cicilanTertunggak": number
    },
    {
      "user_id": "string",
      "nama": "string",
      "cicilanTertunggak": number
    }
  ]
}
```

#### Mengisi Log Tugas

Request :

- Method: `POST`
- Endpoint: `/debt-collector/log-tugas/Create`
- Body :

```json
{
  "tugasId": "string",
  "description": "string"
}
```

Response :

- Status : 200 OK

```json
{
  "responseCode": "string",
  "message": "success"
}
```

#### Get All Log Tugas

Request :

- Method: `GET`
- Endpoint: `/debt-collector/log-tugas`

Response :

- Status : 200 OK
- Body:

```json
{
  "status": "success",
  "log_tugas": [
    {
      "tugasId": "123",
      "description": "Debitur dihubungi, janji akan membayar minggu depan."
    },
    {
      "tugasId": "678",
      "description": "Debitur tidak bisa dihubungi."
    }
    {
      "tugasId": "890",
      "description": "Debitur tidak bisa dihubungi."
    }
  ]
}

```

#### Update Log Tugas

Request :

- Method: `PUT`
- Endpoint: `/debt-collector/log-tugas/:id`
- Body :

```json
{
  "description": "string"
}
```

Response :

- Status : 200 OK

```json
{
  "responseCode": "string",
  "message": "success"
}
```
