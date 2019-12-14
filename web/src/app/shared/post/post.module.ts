import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { StoreModule } from '@ngrx/store';
import { TemplateModule } from '../template/template.module';
import { PostAuthorComponent } from './author';
import { PostCategoriesComponent } from './categories';
import { PostContentComponent } from './content';
import { PostExcerptComponent } from './excerpt';
import { PostFeaturedImageComponent } from './featured-image';
import { PostMetadataComponent } from './metadata';
import * as postReducer from './post.reducer';
import { PostShareToComponent } from './share-to';
import { PostTagsComponent } from './tags';
import { PostTitleComponent } from './title';

@NgModule({
  imports: [
    CommonModule,
    FontAwesomeModule,
    HttpClientModule,
    RouterModule,
    StoreModule.forFeature('post', postReducer.reducer),
    TemplateModule,
  ],
  declarations: [
    PostAuthorComponent,
    PostCategoriesComponent,
    PostContentComponent,
    PostExcerptComponent,
    PostFeaturedImageComponent,
    PostMetadataComponent,
    PostShareToComponent,
    PostTagsComponent,
    PostTitleComponent,
  ],
  exports: [
    PostAuthorComponent,
    PostCategoriesComponent,
    PostContentComponent,
    PostExcerptComponent,
    PostFeaturedImageComponent,
    PostMetadataComponent,
    PostShareToComponent,
    PostTagsComponent,
    PostTitleComponent,
  ],
})
export class PostModule { }
