import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ArchiveComponent } from './archive.component';

const routes: Routes = [
  { path: 'category/:slug', component: ArchiveComponent, data: { type: 'category' } },
  { path: 'tag/:slug', component: ArchiveComponent, data: { type: 'tag' } }, ,
];

@NgModule({
  imports: [
    RouterModule.forChild(routes),
  ],
  exports: [
    RouterModule,
  ],
})
export class ArchiveRoutingModule { }
