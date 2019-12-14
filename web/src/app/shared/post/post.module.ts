import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { StoreModule } from '@ngrx/store';
import { TemplateModule } from '../template/template.module';
import { PostAuthorComponent } from './author';
import { PostMetadataComponent } from './metadata';
import { PostCategoriesComponent } from './post-categories.component';
import { PostContentComponent } from './post-content.component';
import { PostExcerptComponent } from './post-excerpt.component';
import { PostFeaturedImageComponent } from './post-featured-image.component';
import { PostTagsComponent } from './post-tags.component';
import { PostTitleComponent } from './post-title.component';
import * as postReducer from './post.reducer';
import { PostShareToComponent } from './share-to';

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
