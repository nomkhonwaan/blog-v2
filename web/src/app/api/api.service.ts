import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { ApiModule } from './api.module';

@Injectable({
  providedIn: ApiModule,
  deps: [HttpClient],
})
export class ApiService {

  constructor(private http: HttpClient) { }

  /**
   * Call to the RESTful API for deleting uploaded file
   *
   * @param slug string
   */
  deleteFile(slug: string): Observable<any> {
    return this.http.delete<any>(`${environment.url}/api/v2.1/storage/${slug}/delete`);
  }

  /**
   * Call to the RESTful API for uploading file
   *
   * @param file File
   */
  uploadFile(file: File): Observable<Attachment> {
    const formData = new FormData();
    formData.set('file', file, file.name);

    return this.http.post<Attachment>(`${environment.url}/api/v2.1/storage/upload`, formData);
  }

}
